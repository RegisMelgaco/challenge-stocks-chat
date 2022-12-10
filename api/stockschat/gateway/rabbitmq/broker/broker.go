package broker

import (
	"local/challengestockschat/stockschat/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/regismelgaco/go-sdks/erring"
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
