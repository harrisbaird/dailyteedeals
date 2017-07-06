package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
	minio "github.com/minio/minio-go"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

var config = parseConfig()

type Config struct {
	PostgresHost              string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort              int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser              string `env:"POSTGRES_USER,required"`
	PostgresPassword          string `env:"POSTGRES_PASSWORD,required"`
	PostgresDatabase          string `env:"POSTGRES_DATABASE" envDefault:"dailyteedeals"`
	PostgresConnectionTimeout int    `env:"POSTGRES_CONNECTION_TIMEOUT" envDefault:"30"`
	PostgresSSLMode           string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`

	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	AWSS3Bucket        string `env:"AWS_S3_BUCKET,required"`
	AWSS3Endpoint      string `env:"AWS_S3_ENDPOINT" envDefault:"s3.amazonaws.com"`

	ScrapydHost string `env:"SCRAPYD_HOST" envDefault:"scrapyd"`
	ScrapydPort int    `env:"SCRAPYD_PORT" envDefault:"6900"`

	RedisHost string `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort int    `env:"REDIS_PORT" envDefault:"6379"`

	IsProduction bool `env:"PRODUCTION"`
}

func IsProduction() bool {
	return config.IsProduction
}

func IsDevelopment() bool {
	return !IsProduction() || IsTest()
}

func IsTest() bool {
	return flag.Lookup("test.v") != nil
}

func ModeString() string {
	if IsProduction() {
		return "production"
	}

	return "development"
}

func ScrapydURL() string {
	return fmt.Sprintf("http://%s:%d", config.ScrapydHost, config.ScrapydPort)
}

func RedisConnectionString() string {
	return fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
}

func PostgresConnectionString() string {
	return drivers.PostgresBuildQueryString(
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDatabase,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresSSLMode)
}

func NewMinioClient() *minio.Client {
	client, err := minio.New(config.AWSS3Endpoint, config.AWSAccessKeyID, config.AWSSecretAccessKey, true)
	if err != nil {
		panic(err)
	}
	return client
}

func S3Bucket() string {
	return config.AWSS3Bucket
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
