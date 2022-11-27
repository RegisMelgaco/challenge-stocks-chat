package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
)

type Repository interface {
	InsertMessage(context.Context, *chat.Message) error
}
