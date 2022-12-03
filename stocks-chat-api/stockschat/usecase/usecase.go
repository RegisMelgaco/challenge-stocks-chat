package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"sync"

	"go.uber.org/zap"
)

type Usecase interface {
	CreateMessage(context.Context, *entity.Message) error
	Listen(context.Context, entity.Listener) (cleanup func(), err error)
	HandleCreateMessageQueue(l *zap.Logger) error
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
