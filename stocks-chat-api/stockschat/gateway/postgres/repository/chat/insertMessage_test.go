package chat_test

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	chatRepository "local/challengestockschat/stockschat/gateway/postgres/repository/chat"
	"testing"
	"time"

	"github.com/regismelgaco/go-sdks/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Chat_Repository_InsertMessage_Success(t *testing.T) {
	t.Parallel()

	pool := postgres.GetDB(t)

	repo := chatRepository.NewRepository(pool)

	arg := chat.Message{
		Author:  "Douglas Adams",
		Content: "42",
	}

	err := repo.InsertMessage(context.Background(), &arg)

	assert.NoError(t, err)
	if assert.WithinDuration(t, time.Now(), arg.CreatedAt, time.Minute) {
		expected := chat.Message{
			Author:    "Douglas Adams",
			Content:   "42",
			CreatedAt: arg.CreatedAt,
		}

		assert.Equal(t, expected, arg)
	}
}
