package repository

import (
	"context"
	"local/challengestockschat/stockschat/entity"

	"github.com/regismelgaco/go-sdks/erring"
)

func (r repo) ListMessages(ctx context.Context, limit int) ([]entity.Message, error) {
	const sql = `
		select author, created_at, content
		from message
		order by created_at desc
		limit $1;
	`

	rows, err := r.p.Query(ctx, sql, limit)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to select messages from db while listing")
	}

	list := []entity.Message{}
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.Author, &msg.CreatedAt, &msg.Content); err != nil {
			return nil, erring.Wrap(err).Describe("failed to scan message row while listing")
		}

		list = append(list, msg)
	}

	return list, nil
}
