package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"

	"github.com/regismelgaco/go-sdks/erring"
)

func (r repository) InsertMessage(ctx context.Context, msg *chat.Message) error {
	const sql = `
		insert into message (author, content)
		values ($1, $2)
		returning (created_at)
	`
	err := r.p.QueryRow(ctx, sql, msg.Author, msg.Content).Scan(&msg.CreatedAt)
	if err != nil {
		return erring.Wrap(err).
			Describe("failed to write message to database").
			Build()
	}

	return nil
}
