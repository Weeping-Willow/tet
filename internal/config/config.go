package config

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	App              App
	DB               DB
	ExternalServices ExternalServices
}

type App struct {
	PortHTTP string `envconfig:"PORT_HTTP" default:"8080"`
}

type ExternalServices struct {
	EcbRssURL string `envconfig:"ECB_RSS_URL" default:"https://www.bank.lv/vk/ecb_rss.xml"`
}

type DB struct {
	Host         string `envconfig:"DB_HOST" default:""`
	Port         string `envconfig:"DB_PORT" default:""`
	User         string `envconfig:"DB_USER" default:""`
	Password     string `envconfig:"DB_PASSWORD" default:""`
	DatabaseName string `envconfig:"DB_NAME" default:""`
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

func NewDb(cfg Config) (*sqlx.DB, error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.DB.User
	mysqlConfig.Passwd = cfg.DB.Password
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = cfg.DB.Host + ":" + cfg.DB.Port
	mysqlConfig.DBName = cfg.DB.DatabaseName
	mysqlConfig.AllowNativePasswords = true
	mysqlConfig.InterpolateParams = true

	db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "connect to db")
	}

	return db, nil
}
