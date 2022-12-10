package usecase

import (
	"context"
	"strings"

	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
	"local/challengestockschat/stockschat/entity"
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

	logger.FromContext(ctx).Debug("messsage created", zap.Any("msg", *msg))

	u.broker.SendToPublishingMessage(ctx, *msg)

	return nil
}
