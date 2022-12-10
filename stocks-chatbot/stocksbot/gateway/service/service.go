package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"local/stocksbot/stocksbot/entity"
	"local/stocksbot/stocksbot/worker"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/regismelgaco/go-sdks/erring"
	"go.uber.org/zap"
)

type service struct {
	client *retryablehttp.Client
}

func New(l *zap.Logger) worker.StockService {
	client := retryablehttp.NewClient()

	const retryMax = 5 * time.Minute

	client.Backoff = retryablehttp.DefaultBackoff
	client.RetryMax = 5
	client.RetryWaitMax = retryMax
	client.RetryWaitMin = time.Second

	client.Logger = &logger{l}

	return &service{client}
}

func (s service) GetStockEvaluation(ctx context.Context, code string) (entity.StockEvaluation, error) {
	req, err := retryablehttp.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", code),
		nil,
	)
	if err != nil {
		return entity.StockEvaluation{}, erring.Wrap(err).Describe("failed to create request stock evaluation")
	}

	res, err := s.client.Do(req)
	if err != nil {
		return entity.StockEvaluation{}, erring.Wrap(err).Describe("failed to send request")
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with status %v", res.Status)

		return entity.StockEvaluation{}, erring.Wrap(err).With("stock_code", code)
	}

	reader := csv.NewReader(res.Body)

	// distart labels row
	rows, err := reader.ReadAll()
	if err != nil {
		return entity.StockEvaluation{}, erring.Wrap(err).
			Describe("failed while reading evaluation csv rows")
	}

	if len(rows) == 1 && len(rows[0]) == 1 {
		return entity.StockEvaluation{}, erring.Wrap(erring.ErrBadRequest).Describe(rows[0][0])
	}

	if len(rows) != 2 && len(rows[0]) != 8 {
		err := errors.New("stock evaluation response csv with wrong columns count")

		return entity.StockEvaluation{}, erring.Wrap(err)
	}

	r := rows[1]

	return entity.StockEvaluation{
		Code:   r[0],
		Time:   fmt.Sprintf("%s %s", r[1], r[2]),
		Open:   r[3],
		High:   r[4],
		Low:    r[5],
		Close:  r[6],
		Volume: r[7],
	}, nil
}

type logger struct {
	logger *zap.Logger
}

func (l *logger) parseFields(keysAndValues []interface{}) []zap.Field {
	//nolint:gomnd
	pairsCount := len(keysAndValues) / 2

	fields := make([]zap.Field, 0, pairsCount)

	for i := 0; i < pairsCount; i++ {
		key, _ := keysAndValues[i*2].(string)

		fields = append(fields, zap.Any(key, keysAndValues[i*2+1]))
	}

	return fields
}

func (l *logger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, l.parseFields(keysAndValues)...)
}

func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, l.parseFields(keysAndValues)...)
}

func (l *logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, l.parseFields(keysAndValues)...)
}

func (l *logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, l.parseFields(keysAndValues)...)
}
