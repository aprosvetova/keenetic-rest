package keenetic

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"regexp"
)

var rchRe = regexp.MustCompile(`(?m)realm="(.*?)" challenge="(.*?)" session_id="(.*?)" session_cookie="(.*?)"`)

func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func hashSHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func parseAuthorization(h http.Header) (realm, ch, sid, sc string, err error) {
	w := h.Get("WWW-Authenticate")
	if w == "" {
		err = errors.New("wrong header")
		return
	}
	matches := rchRe.FindStringSubmatch(w)
	if len(matches) != 5 {
		err = errors.New("missing realm/challenge")
	}
	realm = matches[1]
	ch = matches[2]
	sid = matches[3]
	sc = matches[4]
	return
}