package chat_test

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	chatRepository "local/challengestockschat/stockschat/gateway/postgres/repository/chat"
	"testing"

	"github.com/regismelgaco/go-sdks/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Repository_Chat_List(t *testing.T) {
	t.Parallel()

	type args struct {
		limit int
	}

	testCases := []struct {
		name string
		args
		seed        []chat.Message
		expected    []chat.Message
		expectedErr error
	}{
		{
			name: "when database is empty expect no error and empty list",
			args: args{
				limit: 10,
			},
			seed:        []chat.Message{},
			expected:    []chat.Message{},
			expectedErr: nil,
		},
		{
			name: "when range is smaller than entries count on messages table expect messages slice with length equals to limit",
			args: args{
				limit: 1,
			},
			seed: []chat.Message{
				{
					Author:    "sheldon",
					Content:   "howdy",
					CreatedAt: time1,
				},
				{
					Author:    "sheldon",
					Content:   "bye",
					CreatedAt: time2,
				},
			},
			expected: []chat.Message{
				{
					Author:    "sheldon",
					Content:   "bye",
					CreatedAt: time2,
				},
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := postgres.GetDB(t)

			if tc.seed != nil {
				SeedMessageTable(t, pool, tc.seed)
			}

			repo := chatRepository.NewRepository(pool)

			actual, err := repo.ListMessages(context.Background(), tc.args.limit)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
