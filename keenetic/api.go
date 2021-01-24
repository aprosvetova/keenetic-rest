package keenetic

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func New(baseURL, login, password string) (*Keenetic, error) {
	k := &Keenetic{
		login:    login,
		password: password,
	}

	jar, _ := cookiejar.New(nil)
	k.c = resty.New()
	k.c.SetHostURL(baseURL).
		SetCookieJar(jar).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(1).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if strings.Contains(r.Request.URL, "/auth") {
				return false
			}

			if r.StatusCode() == http.StatusUnauthorized {
				err := k.auth(r)
				if err == nil {
					return true
				}
			}
			return false
		})

	err := k.auth(nil)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Keenetic) auth(r *resty.Response) (err error) {
	if r == nil {
		r, err = k.c.R().Get("/auth")
		if err != nil {
			return
		}
	}
	realm, ch, sid, sc, err := parseAuthorization(r.Header())
	if err != nil {
		return
	}
	pwd := hashSHA256(ch + hashMD5(fmt.Sprintf("%s:%s:%s", k.login, realm, k.password)))
	k.c.SetCookie(&http.Cookie{
		Name:  sc,
		Value: sid,
		Path:  "/",
	})

	r, err = k.c.R().
		SetBody(map[string]interface{}{"login": k.login, "password": pwd}).
		Post("/auth")
	if err != nil {
		return err
	}
	if r.StatusCode() != http.StatusOK {
		return errors.New("auth error")
	}

	return nil
}

type Keenetic struct {
	c        *resty.Client
	login    string
	password string
}
