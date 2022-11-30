package chat

import (
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
