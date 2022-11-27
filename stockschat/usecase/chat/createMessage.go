package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
)

func (u usecase) CreateMessage(ctx context.Context, msg *chat.Message) error {
	return u.r.InsertMessage(ctx, msg)
}
