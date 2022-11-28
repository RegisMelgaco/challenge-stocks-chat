package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	"sync"
)

type listeners struct {
	mux  sync.Mutex
	pool map[*chat.Listener]interface{}
}

func (l *listeners) add(lis *chat.Listener) {
	l.mux.Lock()
	l.pool[lis] = struct{}{}
	l.mux.Unlock()
}

func (l *listeners) rm(lis *chat.Listener) {
	l.mux.Lock()
	delete(l.pool, lis)
	l.mux.Unlock()
}

func (l *listeners) send(msg chat.Message) {
	for lis := range l.pool {
		f := *lis
		err := f(msg)
		if err != nil {
			l.rm(lis)
		}
	}
}

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
