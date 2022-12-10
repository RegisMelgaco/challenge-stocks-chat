package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"sync"
	"time"

	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

type Usecase interface {
	CreateMessage(context.Context, *entity.Message) error
	Listen(context.Context, entity.Listener) (cleanup func(), err error)
	HandleCreateMessageQueue(l *zap.Logger) error
	HandlePublishMessageQueue(l *zap.Logger) error
}

type usecase struct {
	repo   Repository
	broker Broker
	*listeners
}

func New(repo Repository, broker Broker) Usecase {
	ls := listeners{
		mux:  sync.Mutex{},
		pool: make(map[*entity.Listener]interface{}),
	}

	u := usecase{repo, broker, &ls}

	return u
}

const timeout = time.Minute

func (u usecase) HandlePublishMessageQueue(l *zap.Logger) error {
	return u.broker.ConsumePublishingMessages(l, func(msg entity.Message) error {
		ctx := logger.AddToCtx(context.Background(), l)
		ctx, cancel := context.WithTimeout(ctx, handleCreateMessageQueueTimeout)
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
