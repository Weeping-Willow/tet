package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	App              App
	ExternalServices ExternalServices
}

type App struct {
	PortHTTP string `envconfig:"PORT_HTTP" default:"8080"`
}

type ExternalServices struct {
	EcbRssURL string `envconfig:"ECB_RSS_URL" default:"https://www.bank.lv/vk/ecb_rss.xml"`
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
