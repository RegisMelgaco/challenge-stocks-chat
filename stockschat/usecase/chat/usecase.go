package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
)

type Usecase interface {
	CreateMessage(context.Context, *chat.Message) error
}

type usecase struct {
	r Repository
}

func NewUsecase(r Repository) Usecase {
	return usecase{r}
}
