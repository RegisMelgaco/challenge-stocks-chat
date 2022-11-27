package main

import (
	"context"
	"fmt"
	"local/challengestockschat/stockschat/config"
	stocksChatHTTP "local/challengestockschat/stockschat/gateway/http"
	"local/challengestockschat/stockschat/gateway/postgres/migration"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	authMigrate "github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/migrate"
	"github.com/regismelgaco/go-sdks/erring"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed to load env configs").Build()
		fmt.Println(err)

		return
	}

	var logger *zap.Logger
	if cfg.IsDev {
		lCfg := zap.NewDevelopmentConfig()
		lCfg.DisableStacktrace = true
		lCfg.DisableCaller = true

		logger, err = lCfg.Build()

		erring.SplitStackFromLogs()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		err = erring.Wrap(err).Describe("failed to load env configs").Build()
		fmt.Println(err)

		return
	}

	pgPool, err := pgxpool.Connect(context.Background(), cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to connect to database").
			Err().
			Log(logger, zap.PanicLevel)
	}

	err = authMigrate.Migrate(cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to run auth module migrations").
			Err().
			Log(logger, zap.PanicLevel)
	}

	err = migration.Migrate(cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to run migrations").
			Err().
			Log(logger, zap.PanicLevel)
	}

	router := stocksChatHTTP.SetupRouter(logger, pgPool, cfg)

	logger.Info("server listening")
	if err = http.ListenAndServe(cfg.Host, router); err != nil {
		erring.Wrap(err).
			Describe("server exit with error").
			Err().
			Log(logger, zap.PanicLevel)
	}
}
