package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

type Worker interface {
	ProcessStocksRequest(logger *zap.Logger)
}

type worker struct {
	broker  Broker
	service StockService
}

func New(broker Broker, service StockService) Worker {
	return worker{broker, service}
}

func (w worker) ProcessStocksRequest(l *zap.Logger) {
	ctx := logger.AddToCtx(context.Background(), l)

	for {
		err := w.processOneStockRequest(ctx)
		if err != nil {
			logger := logger.FromContext(ctx)

			_ = erring.Wrap(err).Log(logger, zap.ErrorLevel)
		}
	}
}

const timeout = 15 * time.Minute

func (w worker) processOneStockRequest(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return w.broker.ConsumeStocksRequests(ctx, func(command string) error {
		if !strings.HasPrefix(command, "stock=") {
			return w.broker.CreateMessage(ctx, "unknown command")
		}

		code := command[6:]

		evaluation, err := w.service.GetStockEvaluation(ctx, code)
		if err != nil {
			return w.broker.CreateMessage(ctx, fmt.Sprintf("failed to request %s", code))
		}

		msg := fmt.Sprintf("code: %s | date: %s | high: %s | low: %s | volume: %v | open: %s | close: %s",
			evaluation.Code,
			evaluation.Time,
			evaluation.High,
			evaluation.Low,
			evaluation.Volume,
			evaluation.Open,
			evaluation.Close,
		)

		return w.broker.CreateMessage(ctx, msg)
	})
}
