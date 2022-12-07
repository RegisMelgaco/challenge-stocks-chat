package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"

	"go.uber.org/zap"
)

type Broker interface {
	SendToRequestingBot(ctx context.Context, command string) error
	SendToPublishingMessage(context.Context, entity.Message) error

	ConsumeCreatingMessage(logger *zap.Logger, handler func(msg string) error) error
	ConsumePublishingMessages(logger *zap.Logger, handler func(msg entity.Message) error) error
}
