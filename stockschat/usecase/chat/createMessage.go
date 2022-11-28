package chat

import (
	"context"
	"fmt"
	"local/challengestockschat/stockschat/entity/chat"
)

func (u usecase) CreateMessage(ctx context.Context, msg *chat.Message) error {
	err := u.repo.InsertMessage(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println(len(u.listeners.pool))

	u.listeners.send(*msg)

	return nil
}
