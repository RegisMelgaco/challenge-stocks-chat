package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	"sync"
)

type Usecase interface {
	CreateMessage(context.Context, *chat.Message) error
	Listen(context.Context, chat.Listener) (cleanup func(), err error)
}

type usecase struct {
	repo   Repository
	broker Broker
	*listeners
}

func NewUsecase(repo Repository, broker Broker) Usecase {
	ls := listeners{
		mux:  sync.Mutex{},
		pool: make(map[*chat.Listener]interface{}),
	}

	u := usecase{repo, broker, &ls}

	return u
}
