package repository

import (
	"context"
	"local/challengestockschat/stockschat/entity"

	"github.com/regismelgaco/go-sdks/erring"
)

func (r repo) InsertMessage(ctx context.Context, msg *entity.Message) error {
	const sql = `
		insert into message (author, content)
		values ($1, $2)
		returning (created_at)
	`

	err := r.p.QueryRow(ctx, sql, msg.Author, msg.Content).Scan(&msg.CreatedAt)
	if err != nil {
		return erring.Wrap(err).Describe("failed to write message to database")
	}

	return nil
}
