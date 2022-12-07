package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"strings"
)

func (u usecase) CreateMessage(ctx context.Context, msg *entity.Message) error {
	if strings.HasPrefix(msg.Content, "/") {
		err := u.broker.SendToRequestingBot(ctx, msg.Content[1:])
		if err != nil {
			return err
		}
	}

	err := u.repo.InsertMessage(ctx, msg)
	if err != nil {
		return err
	}

	u.broker.SendToPublishingMessage(ctx, *msg)

	return nil
}
