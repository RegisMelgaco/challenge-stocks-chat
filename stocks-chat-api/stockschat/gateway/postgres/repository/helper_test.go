package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"local/challengestockschat/stockschat/entity"
)

func SeedMessageTable(t *testing.T, pool *pgxpool.Pool, msgs []entity.Message) {
	const sql = `
		insert into message (author, content, created_at)
		values ($1, $2, $3)
	`
	for _, msg := range msgs {
		_, err := pool.Exec(context.Background(), sql, msg.Author, msg.Content, msg.CreatedAt)
		require.NoError(t, err, "failed to seed db with messages")
	}
}

var (
	time1 = time.Time{}.Add(time.Minute)
	time2 = time1.Add(time.Minute)
)
