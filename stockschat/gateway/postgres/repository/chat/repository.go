package chat

import (
	"local/challengestockschat/stockschat/usecase/chat"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	p *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) chat.Repository {
	return repository{pool}
}
