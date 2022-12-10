package broker

import (
	"encoding/json"
	"local/challengestockschat/stockschat/entity"

	"github.com/regismelgaco/go-sdks/erring"
	"go.uber.org/zap"
)

func (b broker) ConsumePublishingMessages(logger *zap.Logger, handler func(msg entity.Message) error) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(b.publishingMessage.Name, "", false, false, false, false, nil)
	if err != nil {
		return erring.Wrap(err).Describe("failed to consume publishingMessage queue")
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			_ = erring.Wrap(err).
				Describe("failed to close channel while ending consume of publishing messages queue").
				Log(logger, zap.ErrorLevel)
		}
	}()

	for d := range msgs {
		var msg publishingMessage
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			return erring.Wrap(err).Describe("failed to unmarshal message from publishing messages queue")
		}

		err := handler(entity.Message{
			Author:    msg.Author,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
		if err != nil {
			nackErr := d.Nack(false, true)
			if nackErr != nil {
				_ = erring.Wrap(nackErr).
					Describe("failed to nack msg from publishing messages queue").
					Log(logger, zap.ErrorLevel)
			}

			return erring.Wrap(err).
				Describe("failed to publish message").
				With("msg", msg)
		} else {
			err := d.Ack(false)
			if err != nil {
				return erring.Wrap(err).Describe("failed to ack publish message")
			}
		}
	}

	return nil
}
