package usecase

import (
	"local/challengestockschat/stockschat/entity"
	"time"

	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const handleCreateMessageQueueTimeout = 5 * time.Minute

func (u usecase) HandleCreateMessageQueue(l *zap.Logger) error {
	return u.broker.ConsumeCreatingMessage(l, func(content string) error {
		ctx := logger.AddToCtx(context.Background(), l)
		ctx, cancel := context.WithTimeout(ctx, handleCreateMessageQueueTimeout)
		defer cancel()

		msg := entity.Message{
			Author:  "bot",
			Content: content,
		}

		return u.CreateMessage(ctx, &msg)
	})
}
