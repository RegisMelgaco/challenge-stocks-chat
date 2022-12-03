package broker

import (
	"local/challengestockschat/stockschat/usecase/chat"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
)

type broker struct {
	conn        *amqp.Connection
	stocksQueue amqp.Queue
}

func New(conn *amqp.Connection) (chat.Broker, error) {
	b := broker{conn: conn}

	ch, err := b.getChannel()
	if err != nil {
		return nil, err
	}

	b.stocksQueue, err = ch.QueueDeclare("requested_stock", true, false, false, false, nil)

	return b, nil
}

func (b broker) getChannel() (*amqp.Channel, error) {
	ch, err := b.conn.Channel()
	if err != nil {
		return nil, erring.Wrap(err).Describe("failed to open a new channel")
	}

	return ch, nil
}
