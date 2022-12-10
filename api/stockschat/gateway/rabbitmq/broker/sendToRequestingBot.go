package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func (b broker) SendToRequestingBot(ctx context.Context, command string) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			logger := logger.FromContext(ctx)
			_ = erring.Wrap(err).
				Describe("failed to close channel of bot command request").
				Log(logger, zap.ErrorLevel)
		}
	}()

	p := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(command),
	}

	err = ch.PublishWithContext(ctx, "", b.requestingBot.Name, false, false, p)
	if err != nil {
		return erring.Wrap(err).Describe("failed to publish msg to stocks queue")
	}

	return nil
}
