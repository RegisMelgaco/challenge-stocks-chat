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
	repo Repository
	*listeners
}

func NewUsecase(repo Repository) Usecase {
	ls := listeners{
		mux:  sync.Mutex{},
		pool: make(map[*chat.Listener]interface{}),
	}

	u := usecase{repo, &ls}

	return u
}
