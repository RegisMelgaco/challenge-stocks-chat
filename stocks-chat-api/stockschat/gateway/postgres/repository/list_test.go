package repository_test

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"local/challengestockschat/stockschat/gateway/postgres/repository"
	"testing"

	"github.com/regismelgaco/go-sdks/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Repository_Chat_List(t *testing.T) {
	t.Parallel()

	for _, tc := range testCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := postgres.GetDB(t)

			if tc.seed != nil {
				SeedMessageTable(t, pool, tc.seed)
			}

			repo := repository.New(pool)

			actual, err := repo.ListMessages(context.Background(), tc.args.limit)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

type args struct {
	limit int
}

func testCases() []struct {
	name string
	args
	seed        []entity.Message
	expected    []entity.Message
	expectedErr error
} {
	return []struct {
		name string
		args
		seed        []entity.Message
		expected    []entity.Message
		expectedErr error
	}{
		{
			name: "when database is empty expect no error and empty list",
			args: args{
				limit: 10,
			},
			seed:        []entity.Message{},
			expected:    []entity.Message{},
			expectedErr: nil,
		},
		{
			name: "when range is smaller than entries count on messages table expect messages slice with length equals to limit",
			args: args{
				limit: 1,
			},
			seed: []entity.Message{
				{
					Author:    "sheldon",
					Content:   "howdy",
					CreatedAt: time1(),
				},
				{
					Author:    "sheldon",
					Content:   "bye",
					CreatedAt: time2(),
				},
			},
			expected: []entity.Message{
				{
					Author:    "sheldon",
					Content:   "bye",
					CreatedAt: time2(),
				},
			},
			expectedErr: nil,
		},
	}
}
