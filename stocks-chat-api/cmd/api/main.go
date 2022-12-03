package main

import (
	"context"
	"local/challengestockschat/stockschat/config"
	stocksChatHTTP "local/challengestockschat/stockschat/gateway/http"
	"local/challengestockschat/stockschat/gateway/postgres/migration"
	"local/challengestockschat/stockschat/gateway/rabbitmq/broker"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	authMigrate "github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/migrate"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed to load env configs")
		log.Println(err.Error())

		return
	}

	logger, err := logger.New(cfg.IsDev)
	if err != nil {
		log.Println(err.Error())

		return
	}

	if err != nil {
		_ = erring.Wrap(err).
			Describe("failed to load env configs").
			Log(logger, zap.PanicLevel)
	}

	amqpConn, err := amqp.Dial(cfg.RabbitMQConn)
	if err != nil {
		_ = erring.Wrap(err).
			Describe("failed to connect to rabbitmq").
			Log(logger, zap.PanicLevel)
	}

	broker, err := broker.New(amqpConn)
	if err != nil {
		_ = erring.Wrap(err).Log(logger, zap.PanicLevel)
	}

	pgPool, err := pgxpool.Connect(context.Background(), cfg.PostgresConn)
	if err != nil {
		_ = erring.Wrap(err).
			Describe("failed to connect to database").
			Log(logger, zap.PanicLevel)
	}

	err = authMigrate.Migrate(cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to run auth module migrations").
			Log(logger, zap.PanicLevel)
	}

	err = migration.Migrate(cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to run migrations").
			Log(logger, zap.PanicLevel)
	}

	router := stocksChatHTTP.SetupRouter(logger, pgPool, cfg, broker)

	logger.Info("server listening")
	if err = http.ListenAndServe(cfg.Host, router); err != nil {
		erring.Wrap(err).
			Describe("server exit with error").
			Log(logger, zap.PanicLevel)
	}
}