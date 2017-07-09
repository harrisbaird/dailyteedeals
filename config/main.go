package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
	minio "github.com/minio/minio-go"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

// App contains app configuration
var App = parseConfig()

type Config struct {
	PostgresHost              string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort              int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser              string `env:"POSTGRES_USER,required"`
	PostgresPassword          string `env:"POSTGRES_PASSWORD,required"`
	PostgresDatabase          string `env:"POSTGRES_DATABASE" envDefault:"dailyteedeals"`
	PostgresConnectionTimeout int    `env:"POSTGRES_CONNECTION_TIMEOUT" envDefault:"30"`
	PostgresSSLMode           string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`

	DomainAPI    string `env:"DOMAIN_API" envDefault:"api.dailyteedeals.com"`
	DomainGo     string `env:"DOMAIN_GO" envDefault:"go.dailyteedeals.com"`
	DomainImages string `env:"DOMAIN_IMAGES" envDefault:"images-17e7.kxcdn.com"`

	HTTPListenAddr string `env:"HTTP_LISTEN_ADDR" envDefault:"0.0.0.0:8080"`

	Env string `env:"APP_ENV" envDefault:"development"`

	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	AWSS3Bucket        string `env:"AWS_S3_BUCKET,required"`
	AWSS3Endpoint      string `env:"AWS_S3_ENDPOINT" envDefault:"s3.amazonaws.com"`

	ScrapydHost string `env:"SCRAPYD_HOST" envDefault:"scrapyd"`
	ScrapydPort int    `env:"SCRAPYD_PORT" envDefault:"6900"`

	RedisHost string `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort int    `env:"REDIS_PORT" envDefault:"6379"`
}

func IsProduction() bool {
	return App.Env == "production"
}

func IsDevelopment() bool {
	return !IsProduction() || IsTest()
}

func IsTest() bool {
	return flag.Lookup("test.v") != nil
}

func EnvironmentString() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}

func ScrapydURL() string {
	return fmt.Sprintf("http://%s:%d", App.ScrapydHost, App.ScrapydPort)
}

func RedisConnectionString() string {
	return fmt.Sprintf("%s:%d", App.RedisHost, App.RedisPort)
}

func PostgresConnectionString() string {
	return drivers.PostgresBuildQueryString(
		App.PostgresUser,
		App.PostgresPassword,
		App.PostgresDatabase,
		App.PostgresHost,
		App.PostgresPort,
		App.PostgresSSLMode)
}

func PostgresTestConnectionString() string {
	return drivers.PostgresBuildQueryString(
		"postgres",
		"",
		"dailyteedeals_test",
		"127.0.0.1",
		5432,
		"disable")
}

func NewMinioClient() *minio.Client {
	client, err := minio.New(App.AWSS3Endpoint, App.AWSAccessKeyID, App.AWSSecretAccessKey, true)
	if err != nil {
		panic(err)
	}
	return client
}

func S3Bucket() string {
	return App.AWSS3Bucket
}

func parseConfig() *Config {
	config := Config{}
	err := env.Parse(&config)
	// Don't panic if missing env vars in test mode
	if err != nil && !IsTest() {
		panic(err)
	}
	return &config
}
