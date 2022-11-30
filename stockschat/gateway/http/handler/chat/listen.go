package chat

import (
	"context"
	"local/challengestockschat/stockschat/entity/chat"
	v1 "local/challengestockschat/stockschat/gateway/http/handler/chat/v1"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

// TODO: handle host origin
func (h Handler) Listen(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = erring.Wrap(err).
			Describe("failed to upgrade connection for listen request").
			Build()

		httpresp.Internal(err).Handle(w, r)

		return
	}

	logger := logger.FromContext(r.Context())

	cleanup, err := h.u.Listen(r.Context(), func(m chat.Message) error {
		err := conn.WriteJSON(v1.ToMessangeOutput(m))

		if err != nil {
			errBuilder := erring.Wrap(err).
				Describe("failed to write message")

			errBuilder.Err().Log(logger, zap.ErrorLevel)

			return errBuilder.Build()
		}

		return nil
	})
	if err != nil {
		conn.Close()

		erring.Wrap(err).Err().Log(logger, zap.ErrorLevel)

		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		cleanup()

		message := websocket.FormatCloseMessage(code, "")
		err := conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		if err != nil {
			errBuilder := erring.Wrap(err).
				Describe("failed to write close control message")

			errBuilder.Err().Log(logger, zap.ErrorLevel)

			return errBuilder.Build()
		}

		return nil
	})

	go func() {
		for {
			var input v1.InputMessage
			if err := conn.ReadJSON(&input); err != nil {
				conn.Close()

				erring.Wrap(err).
					Describe("failed to read message").
					Err().Log(logger, zap.ErrorLevel)

				break
			}

			msg, err := input.ToEntity(r.Context())
			if err != nil {
				conn.Close()

				erring.Wrap(err).
					Describe("failed to encode message message as json").
					Err().Log(logger, zap.ErrorLevel)

				break
			}

			if err := h.u.CreateMessage(context.Background(), &msg); err != nil {
				conn.Close()

				erring.Wrap(err).
					Err().Log(logger, zap.ErrorLevel)

				break
			}
		}
	}()
}
