package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"time"

	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

const handlePublishMessageQueueTimeout = time.Minute

func (u usecase) HandlePublishMessageQueue(l *zap.Logger) error {
	return u.broker.ConsumePublishingMessages(l, func(msg entity.Message) error {
		ctx := logger.AddToCtx(context.Background(), l)
		ctx, cancel := context.WithTimeout(ctx, handlePublishMessageQueueTimeout)
		defer cancel()

		for lis := range u.listeners.pool {
			select {
			case <-ctx.Done():
				return erring.Wrap(entity.ErrTimeout)
			default:
				f := *lis
				err := f(msg)
				if err != nil {
					u.listeners.rm(lis)
				}
			}
		}

		return nil
	})
}
