package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
)

func (b broker) RequestBotCommand(ctx context.Context, command string) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	defer ch.Close()

	p := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(command),
	}

	err = ch.PublishWithContext(ctx, "", b.stocksQueue.Name, false, false, p)
	if err != nil {
		return erring.Wrap(err).Describe("failed to publish msg to stocks queue")
	}

	return nil
}
