package broker

import (
	"github.com/regismelgaco/go-sdks/erring"
	"go.uber.org/zap"
)

func (b broker) ConsumeCreatingMessage(logger *zap.Logger, handler func(msg string) error) error {
	ch, err := b.getChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(b.creatingMessage.Name, "", false, false, false, false, nil)
	if err != nil {
		return erring.Wrap(err).Describe("failed to consume creatingMessage queue")
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			_ = erring.Wrap(err).
				Describe("failed to close channel while ending consume of stocks requests").
				Log(logger, zap.ErrorLevel)
		}
	}()

	for d := range msgs {
		msg := string(d.Body)

		err := handler(msg)
		if err != nil {
			nackErr := d.Nack(false, true)
			if nackErr != nil {
				_ = erring.Wrap(nackErr).
					Describe("failed to nack create message request").
					Log(logger, zap.ErrorLevel)
			}

			return erring.Wrap(err).
				Describe("failed to handle create message request").
				With("msg", msg)
		} else {
			err := d.Ack(false)
			if err != nil {
				return erring.Wrap(err).Describe("failed to ack create message request")
			}
		}
	}

	return nil
}
