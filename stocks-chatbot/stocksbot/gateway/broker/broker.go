package broker

import (
	"context"
	"local/stocksbot/stocksbot/worker"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

type broker struct {
	conn                *amqp.Connection
	requestedStockQueue amqp.Queue
	createMessageQueue  amqp.Queue
}

func New(conn *amqp.Connection) (worker.Broker, error) {
	b := &broker{conn: conn}

	ch, err := b.getChannel()
	if err != nil {
		return nil, err
	}

	b.requestedStockQueue, err = ch.QueueDeclare("requested_stock", true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create requested_stock queue")
	}

	b.createMessageQueue, err = ch.QueueDeclare("create_message", true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to create create_message queue")
	}

	return b, nil
}

func (b *broker) ConsumeStocksRequests(ctx context.Context, handler func(string) error) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(b.requestedStockQueue.Name, "", false, false, false, false, nil)
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

	err = ch.Publish("", b.createMessageQueue.Name, false, false, p)
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
