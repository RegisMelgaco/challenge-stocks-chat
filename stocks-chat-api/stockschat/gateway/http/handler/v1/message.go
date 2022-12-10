package v1

import (
	"context"
	"time"

	"github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/erring"
	"local/challengestockschat/stockschat/entity"
)

type InputMessage struct {
	Content string `json:"content"`
}

func (i InputMessage) ToEntity(ctx context.Context) (entity.Message, error) {
	claims, err := handler.ClaimsFromContext(ctx)
	if err != nil {
		return entity.Message{}, erring.Wrap(err).
			Describe("failed to get claims from context")
	}

	return entity.Message{
		Content: i.Content,
		Author:  claims.UserName,
	}, nil
}

type OutputMessage struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func ToMessangeOutput(msg entity.Message) Response {
	return Response{
		Type: "message",
		Payload: OutputMessage{
			Author:    msg.Author,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		},
	}
}
