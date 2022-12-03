package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
)

func (u usecase) CreateMessage(ctx context.Context, msg *entity.Message) error {
	if msg.Content[0] == '/' {
		err := u.broker.RequestBotCommand(ctx, msg.Content[1:])
		if err != nil {
			return err
		}
	}

	err := u.repo.InsertMessage(ctx, msg)
	if err != nil {
		return err
	}

	u.listeners.send(*msg)

	return nil
}
