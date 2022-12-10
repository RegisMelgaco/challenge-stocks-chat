package usecase

import (
	"context"
	"local/challengestockschat/stockschat/entity"
)

const listMessagesMax = 50

func (u usecase) Listen(ctx context.Context, lis entity.Listener) (func(), error) {
	u.listeners.add(&lis)

	cleanup := func() {
		u.listeners.rm(&lis)
	}

	msgs, err := u.repo.ListMessages(ctx, listMessagesMax)
	if err != nil {
		cleanup()

		return nil, err
	}

	for _, msg := range msgs {
		err := lis(msg)
		if err != nil {
			cleanup()

			return nil, err
		}
	}

	return cleanup, nil
}
