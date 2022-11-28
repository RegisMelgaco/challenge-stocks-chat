package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"

	"github.com/regismelgaco/go-sdks/erring"
)

func (r repository) ListMessages(ctx context.Context, limit int) ([]chat.Message, error) {
	const sql = `
		select author, created_at, content
		from message
		order by created_at desc
		limit $1;
	`

	rows, err := r.p.Query(ctx, sql, limit)
	if err != nil {
		return nil, erring.Wrap(err).
			Describe("failed to select messages from db while listing").
			Build()
	}

	list := []chat.Message{}
	for rows.Next() {
		var msg chat.Message
		if err := rows.Scan(&msg.Author, &msg.CreatedAt, &msg.Content); err != nil {
			return nil, erring.Wrap(err).
				Describe("failed to scan message row while listing").
				Build()
		}

		list = append(list, msg)
	}

	return list, nil
}
