package usecase

import (
	"context"

	"local/challengestockschat/stockschat/entity"
)

type Repository interface {
	InsertMessage(context.Context, *entity.Message) error
	ListMessages(ctx context.Context, limit int) ([]entity.Message, error)
}
