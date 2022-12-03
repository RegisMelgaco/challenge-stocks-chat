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
		err = erring.Wrap(err).Describe("failed to upgrade connection for listen request")

		httpresp.Internal(err).Handle(w, r)

		return
	}

	logger := logger.FromContext(r.Context())

	cleanup, err := h.u.Listen(r.Context(), func(m chat.Message) error {
		err := conn.WriteJSON(v1.ToMessangeOutput(m))

		if err != nil {
			return erring.Wrap(err).Log(logger, zap.ErrorLevel)
		}

		return nil
	})
	if err != nil {
		conn.Close()

		_ = erring.Wrap(err).Log(logger, zap.ErrorLevel)

		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		cleanup()

		message := websocket.FormatCloseMessage(code, "")
		err := conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		if err != nil {
			return erring.Wrap(err).
				Describe("failed to write close control message").
				Log(logger, zap.ErrorLevel)
		}

		return nil
	})

	go func() {
		for {
			var input v1.InputMessage
			if err := conn.ReadJSON(&input); err != nil {
				conn.Close()

				_ = erring.Wrap(err).
					Describe("failed to read message").
					Log(logger, zap.ErrorLevel)

				break
			}

			msg, err := input.ToEntity(r.Context())
			if err != nil {
				conn.Close()

				_ = erring.Wrap(err).
					Describe("failed to encode message message as json").
					Log(logger, zap.ErrorLevel)

				break
			}

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
			defer cancel()

			if err := h.u.CreateMessage(ctx, &msg); err != nil {
				conn.Close()

				_ = erring.Wrap(err).Log(logger, zap.ErrorLevel)

				break
			}
		}
	}()
}
