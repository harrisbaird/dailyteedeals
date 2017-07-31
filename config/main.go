package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/go-pg/pg"
)

// App contains app configuration
var App = parseConfig()

type Config struct {
	Env string `env:"APP_ENV" envDefault:"development"`

	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:""`
	PostgresDatabase string `env:"POSTGRES_DATABASE" envDefault:"dailyteedeals"`

	DomainWeb    string `env:"DOMAIN_WEB" envDefault:"dailyteedeals.com"`
	DomainAPI    string `env:"DOMAIN_API" envDefault:"api.dailyteedeals.com"`
	DomainGo     string `env:"DOMAIN_GO" envDefault:"go.dailyteedeals.com"`
	DomainImages string `env:"DOMAIN_IMAGES" envDefault:"cdn.dailyteedeals.com"`

	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	AWSS3Bucket        string `env:"AWS_S3_BUCKET" envDefault:"dailyteedeals"`
	AWSS3Endpoint      string `env:"AWS_S3_ENDPOINT" envDefault:"s3.amazonaws.com"`
	AWSS3Secure        bool   `env:"AWS_S3_SECURE" envDefault:"true"`

	HTTPListenAddr string `env:"HTTP_LISTEN_ADDR" envDefault:"0.0.0.0:8080"`
	ScrapydAddr    string `env:"SCRAPYD_ADDR" envDefault:"scrapyd:6900"`
	RedisAddr      string `env:"REDIS_ADDR" envDefault:"redis:6379"`

	ItemsPerPage int `env:"ITEMS_PER_PAGE" envDefault:"200"`
}

func IsProduction() bool {
	return App.Env == "production"
}

func PostgresConnectionOptions() *pg.Options {
	return &pg.Options{
		Addr:     fmt.Sprintf("%s:%d", App.PostgresHost, App.PostgresPort),
		User:     App.PostgresUser,
		Password: App.PostgresPassword,
		Database: App.PostgresDatabase,
	}

}

func PostgresTestConnectionOptions() *pg.Options {
	return &pg.Options{
		Addr:     "0.0.0.0:5432",
		User:     "postgres",
		Password: "",
		Database: "dailyteedeals_test",
		PoolSize: 20,
	}
}

func parseConfig() *Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil && config.Env == "production" {
		panic(err)
	}
	return &config
}
