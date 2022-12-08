package main

import (
	"fmt"
	"local/stocksbot/stocksbot/config"
	"local/stocksbot/stocksbot/gateway/broker"
	"local/stocksbot/stocksbot/gateway/service"
	"local/stocksbot/stocksbot/worker"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		err = erring.Wrap(err).Describe("failed to load env configs")
		fmt.Println(err)

		return
	}

	logger, err := logger.New(cfg.IsDev)
	if err != nil {
		panic(err)
	}

	if err != nil {
		_ = erring.Wrap(err).
			Describe("failed to create logger").
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

	service := service.New(logger)

	worker := worker.New(broker, service)

	logger.Info("worker starting")

	worker.ProcessStocksRequest(logger)
}
