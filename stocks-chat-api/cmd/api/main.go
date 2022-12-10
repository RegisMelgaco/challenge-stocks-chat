package main

import (
	"context"
	"local/challengestockschat/stockschat/config"
	"local/challengestockschat/stockschat/gateway/postgres/migration"
	"local/challengestockschat/stockschat/gateway/postgres/repository"
	"local/challengestockschat/stockschat/gateway/rabbitmq/broker"
	"local/challengestockschat/stockschat/usecase"
	"log"
	"net/http"

	stocksChatHTTP "local/challengestockschat/stockschat/gateway/http"

	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
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

	if err != nil {
		log.Panicln("failed to load env configs", err)
	}

	logger, err := logger.New(cfg.IsDev)
	if err != nil {
		log.Println(err.Error())

		return
	}

	start(logger, cfg)
}

func start(logger *zap.Logger, cfg config.Config) {
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

	err = migration.Migrate(cfg.PostgresConn)
	if err != nil {
		erring.Wrap(err).
			Describe("failed to run migrations").
			Log(logger, zap.PanicLevel)
	}

	usecase := usecase.New(repository.New(pgPool), broker)

	go func() {
		logger.Info("started consuming from create messages queue")

		if err := usecase.HandleCreateMessageQueue(logger); err != nil {
			erring.Wrap(err).
				Describe("stopped handling create messages request from queue").
				Log(logger, zap.PanicLevel)
		}
	}()

	go func() {
		logger.Info("started consuming from publish messages queue")

		if err := usecase.HandlePublishMessageQueue(logger); err != nil {
			erring.Wrap(err).
				Describe("stopped handling publish messages request from queue").
				Log(logger, zap.PanicLevel)
		}
	}()

	router := stocksChatHTTP.SetupRouter(logger, pgPool, cfg, usecase)

	logger.Info("server listening")

	if err := http.ListenAndServe(cfg.Host, router); err != nil {
		erring.Wrap(err).
			Describe("server exit with error").
			Log(logger, zap.PanicLevel)
	}
}
