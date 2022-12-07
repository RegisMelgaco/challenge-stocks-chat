package broker

import (
	"context"
	"local/stocksbot/stocksbot/worker"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

const (
	requestingBotQueue   = "requesting-bot"
	creatingMessageQueue = "creating-message"
)

type broker struct {
	conn            *amqp.Connection
	requestingBot   amqp.Queue
	creatingMessage amqp.Queue
}

func New(conn *amqp.Connection) (worker.Broker, error) {
	b := &broker{conn: conn}

	ch, err := b.getChannel()
	if err != nil {
		return nil, err
	}

	b.requestingBot, err = ch.QueueDeclare(requestingBotQueue, true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create requesting-bot queue")
	}

	b.creatingMessage, err = ch.QueueDeclare(creatingMessageQueue, true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create creating-message queue")
	}

	return b, nil
}

func (b *broker) ConsumeStocksRequests(ctx context.Context, handler func(string) error) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(b.requestingBot.Name, "", false, false, false, false, nil)
	defer func() {
		err := ch.Close()
		if err != nil {
			logger := logger.FromContext(ctx)
			_ = erring.Wrap(err).
				Describe("failed to close channel while ending consume of stocks requests").
				Log(logger, zap.ErrorLevel)
		}
	}()

	for d := range msgs {
		code := string(d.Body)

		err := handler(code)
		if err != nil {
			nackErr := d.Nack(false, true)
			if nackErr != nil {
				logger := logger.FromContext(ctx)

				_ = erring.Wrap(nackErr).
					Describe("failed to nack stock request").
					Log(logger, zap.ErrorLevel)
			}

			return erring.Wrap(err).
				Describe("failed to handle stock request").
				With("stock_code", code)
		} else {
			err := d.Ack(false)
			if err != nil {
				return erring.Wrap(err).Describe("failed to ack stock request")
			}
		}
	}

	return nil
}

func (b *broker) CreateMessage(ctx context.Context, msg string) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	p := amqp.Publishing{
		ContentType:  "plain/text",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(msg),
	}

	err = ch.Publish("", b.creatingMessage.Name, false, false, p)
	if err != nil {
		return erring.Wrap(err).Describe("failed to deliver create message request")
	}

	return nil
}

func (b *broker) getChannel() (*amqp.Channel, error) {
	ch, err := b.conn.Channel()
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to get amqp channel")
	}

	return ch, nil
}
