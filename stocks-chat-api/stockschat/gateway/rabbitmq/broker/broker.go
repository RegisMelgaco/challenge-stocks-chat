package broker

import (
	"context"
	"encoding/json"
	"local/challengestockschat/stockschat/entity"
	"local/challengestockschat/stockschat/usecase"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

type broker struct {
	conn *amqp.Connection

	requestingBot     amqp.Queue
	creatingMessage   amqp.Queue
	publishingMessage amqp.Queue
}

const (
	requestingBotQueue        = "requesting-bot"
	creatingMessageQueue      = "creating-message"
	publishingMessageExchange = "publishing-message"
)

func New(conn *amqp.Connection) (usecase.Broker, error) {
	b := broker{conn: conn}

	ch, err := b.getChannel()
	if err != nil {
		return nil, err
	}

	b.requestingBot, err = ch.QueueDeclare(requestingBotQueue, true, false, false, false, nil)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create requesting-bot queue")
	}

	b.creatingMessage, err = ch.QueueDeclare(creatingMessageQueue, true, false, false, false, nil)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create creating-message queue")
	}

	err = ch.ExchangeDeclare(publishingMessageExchange, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create publishing-message exchange")
	}

	b.publishingMessage, err = ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create publishing-message queue")
	}

	err = ch.QueueBind(b.publishingMessage.Name, "", publishingMessageExchange, false, nil)
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to bind publishing message queue to exchange")
	}

	return b, nil
}

func (b broker) getChannel() (*amqp.Channel, error) {
	ch, err := b.conn.Channel()
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to open a new channel")
	}

	return ch, nil
}

type publishingMessage struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"create_at"`
}

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

	p := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         bytes,
	}

	err = ch.PublishWithContext(ctx, publishingMessageExchange, "", false, false, p)
	if err != nil {
		return erring.Wrap(err).Describe("failed to publish msg to publishing messages queue")
	}

	return nil
}

func (b broker) ConsumePublishingMessages(logger *zap.Logger, handler func(msg entity.Message) error) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(b.publishingMessage.Name, "", false, false, false, false, nil)
	defer func() {
		err := ch.Close()
		if err != nil {
			_ = erring.Wrap(err).
				Describe("failed to close channel while ending consume of publishing messages queue").
				Log(logger, zap.ErrorLevel)
		}
	}()

	for d := range msgs {
		var msg entity.Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			return erring.Wrap(err).Describe("failed to unmarshal message from publishing messages queue")
		}

		err := handler(msg)
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
