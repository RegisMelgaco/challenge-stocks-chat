package v1

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	"time"

	"github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/erring"
)

type InputMessage struct {
	Content string `json:"content"`
}

func (i InputMessage) ToEntity(ctx context.Context) (chat.Message, error) {
	claims, err := handler.ClaimsFromContext(ctx)
	if err != nil {
		return chat.Message{}, erring.Wrap(err).
			Describe("failed to get claims from context")
	}

	return chat.Message{
		Content: i.Content,
		Author:  claims.UserName,
	}, nil
}

type OutputMessage struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func ToMessangeOutput(msg chat.Message) OutputMessage {
	return OutputMessage{
		Author:    msg.Author,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}
}
