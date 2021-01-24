package main

import (
	"fmt"
	"github.com/aprosvetova/keenetic-rest/keenetic"
	"github.com/caarlos0/env/v6"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var cfg config
var k *keenetic.Keenetic

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	var err error

	k, err = keenetic.New(cfg.BaseURL, cfg.Login, cfg.Password)
	if err != nil {
		log.Fatalln("Keenetic", err)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func isGuestOn() bool {
	ifaces, err := k.GetInterfaces()
	if err != nil {
		return false
	}
	for i, up := range ifaces {
		for _, ic := range cfg.Interfaces {
			if i == ic && up {
				return true
			}
		}
	}
	return false
}

func setGuest(up bool) {
	ifaces := make(map[string]bool)
	for _, i := range cfg.Interfaces {
		ifaces[i] = up
	}
	k.SetInterfaces(ifaces)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, strconv.FormatBool(isGuestOn()))
		return
	}
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, "false")
			return
		}
		bod := string(body)
		if bod == "true" {
			setGuest(true)
			fmt.Fprint(w, "true")
			return
		} else if bod == "false" {
			setGuest(false)
			fmt.Fprint(w, "false")
			return
		}
		fmt.Fprint(w, "false")
		return
	}
}