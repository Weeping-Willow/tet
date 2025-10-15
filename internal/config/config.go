package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	App App
}

type App struct {
	PortHTTP string `envconfig:"PORT_HTTP" default:"8080"`
}

func New() (Config, error) {
	cfg := Config{}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("no .env file found")
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "cannot process config")
	}

	return cfg, nil
}
