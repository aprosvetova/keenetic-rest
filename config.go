package main

type config struct {
	BaseURL string `env:"ROUTER_URL,required"`
	Login string `env:"ROUTER_LOGIN,required"`
	Password string `env:"ROUTER_PASSWORD,required"`

	Interfaces []string `env:"INTERFACES,required"`
}
