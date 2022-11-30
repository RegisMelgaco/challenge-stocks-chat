package chat

import "context"

type Broker interface {
	RequestBotCommand(ctx context.Context, command string) error
}
