package main

import (
	"context"
	"fmt"
	"local/challengestockschat/stockschat/config"
	stocksChatHTTP "local/challengestockschat/stockschat/gateway/http"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	authMigrate "github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/migrate"
	"github.com/regismelgaco/go-sdks/erring"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed to load env configs").Build()
		fmt.Println(err)

		return
	}

	pgPool, err := pgxpool.Connect(context.Background(), cfg.PostgresConn)
	if err != nil {
		err = erring.Wrap(err).Describe("failed to connect to database").Build()
		fmt.Println(err)

		return
	}

	err = authMigrate.Migrate(cfg.PostgresConn)
	if err != nil {
		err = erring.Wrap(err).Describe("failed to run auth module migrations").Build()
		fmt.Println(err)

		return
	}

	router := stocksChatHTTP.SetupRouter(pgPool, cfg)

	fmt.Println("starting server")
	if err = http.ListenAndServe(cfg.Host, router); err != nil {
		err = erring.Wrap(err).Describe("server exit with error").Build()
		fmt.Println(err)

		return
	}
}
