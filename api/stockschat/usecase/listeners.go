package usecase

import (
	"local/challengestockschat/stockschat/entity"
	"sync"
)

type listeners struct {
	mux  sync.Mutex
	pool map[*entity.Listener]interface{}
}

func (l *listeners) add(lis *entity.Listener) {
	l.mux.Lock()
	l.pool[lis] = struct{}{}
	l.mux.Unlock()
}

func (l *listeners) rm(lis *entity.Listener) {
	l.mux.Lock()
	delete(l.pool, lis)
	l.mux.Unlock()
}

func (l *listeners) send(msg entity.Message) {
	for lis := range l.pool {
		f := *lis

		err := f(msg)
		if err != nil {
			l.rm(lis)
		}
	}
}
