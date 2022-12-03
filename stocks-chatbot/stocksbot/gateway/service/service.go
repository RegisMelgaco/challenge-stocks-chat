package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"local/stocksbot/stocksbot/entity"
	"local/stocksbot/stocksbot/worker"
	"net/http"

	"github.com/regismelgaco/go-sdks/erring"
)

type service struct {
	client http.Client
}

func New() worker.StockService {
	return &service{
		client: http.Client{},
	}
}

func (s service) GetStockEvaluation(ctx context.Context, code string) (entity.StockEvaluation, error) {
	req, err := http.NewRequest(
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
