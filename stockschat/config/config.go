package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/regismelgaco/go-sdks/erring"
)

type Config struct {
	Host         string `envconfig:"HOST" default:":3000"`
	PostgresConn string `envconfig:"POSTGRES_CONN" required:"true"`
	JWTSecret    string `envconfig:"JWT_SECRET" required:"true"`
	IsDev        bool   `envconfig:"IS_DEV"`
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed while trying to load env vars").Build()

		fmt.Println(err)
	}

	var cfg Config
	err = envconfig.Process(".", &cfg)
	if err != nil {
		return Config{}, erring.Wrap(err).Describe("failed while loading configs").Build()
	}

	return cfg, nil
}
