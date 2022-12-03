package worker

import "context"

type Broker interface {
	ConsumeStocksRequests(ctx context.Context, handler func(command string) error) error
	CreateMessage(ctx context.Context, msg string) error
}
