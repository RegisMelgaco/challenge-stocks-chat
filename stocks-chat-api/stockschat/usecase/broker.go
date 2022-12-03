package usecase

import (
	"context"

	"go.uber.org/zap"
)

type Broker interface {
	RequestBotCommand(ctx context.Context, command string) error
	ConsumeCreateMessage(logger *zap.Logger, handler func(msg string) error) error
}
