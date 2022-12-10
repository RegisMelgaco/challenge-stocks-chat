package broker

import (
	"context"
	"encoding/json"
	"local/challengestockschat/stockschat/entity"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func (b broker) SendToPublishingMessage(ctx context.Context, msg entity.Message) error {
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

	bytes, err := json.Marshal(publishingMessage{
		Author:    msg.Author,
		CreatedAt: msg.CreatedAt,
		Content:   msg.Content,
	})
	if err != nil {
		return erring.Wrap(err).Describe("failed to decode publishing messsage")
	}

	p := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         bytes,
	}

	err = ch.PublishWithContext(ctx, publishingMessageExchange, "", false, false, p)
	if err != nil {
		return erring.Wrap(err).Describe("failed to publish msg to publishing messages queue")
	}

	return nil
}
