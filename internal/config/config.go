package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type ArgsHTTP struct {
	Address  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	URLShort string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

var AFlag ArgsHTTP

func InitializeArgs() error {
	flag.StringVar(&AFlag.Address, "a", AFlag.Address, "Address HTTP-server")
	flag.StringVar(&AFlag.URLShort, "b", AFlag.URLShort, "Base URL")

	if err := env.Parse(&AFlag); err != nil {
		return err
	}
	return nil
}
