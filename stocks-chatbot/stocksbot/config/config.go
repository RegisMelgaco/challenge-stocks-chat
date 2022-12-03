package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/regismelgaco/go-sdks/erring"
)

type Config struct {
	Host         string `envconfig:"HOST" default:":3000"`
	RabbitMQConn string `envconfig:"RABBIT_MQ_CONN" required:"true"`
	IsDev        bool   `envconfig:"IS_DEV"`
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed while trying to load env vars")

		log.Println(err)
	}

	var cfg Config
	err = envconfig.Process(".", &cfg)
	if err != nil {
		return Config{}, erring.Wrap(err).Describe("failed while loading configs")
	}

	return cfg, nil
}
