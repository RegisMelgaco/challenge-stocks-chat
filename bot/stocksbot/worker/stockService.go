package worker

import (
	"context"
	"local/stocksbot/stocksbot/entity"
)

type StockService interface {
	GetStockEvaluation(ctx context.Context, stockCode string) (entity.StockEvaluation, error)
}
