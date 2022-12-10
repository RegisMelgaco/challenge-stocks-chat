package repository

import (
	"local/challengestockschat/stockschat/usecase"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	p *pgxpool.Pool
}

func New(pool *pgxpool.Pool) usecase.Repository {
	return repo{pool}
}
