package chat

import (
	"context"
	"fmt"
	"local/challengestockschat/stockschat/entity/chat"
	v1 "local/challengestockschat/stockschat/gateway/http/handler/chat/v1"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
)

func (h Handler) Listen(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = erring.Wrap(err).
			Describe("failed to upgrade connection for listen request").
			Build()

		httpresp.Internal(err).Handle(w, r)

		return
	}

	cleanup, err := h.u.Listen(r.Context(), func(m chat.Message) error {
		err := conn.WriteJSON(v1.ToMessangeOutput(m))

		if err != nil {
			err = erring.Wrap(err).
				Describe("failed to write message").
				Build()

			return err
		}

		return nil
	})
	if err != nil {
		conn.Close()

		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		cleanup()

		message := websocket.FormatCloseMessage(code, "")
		conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))

		return nil
	})

	go func() {
		for {
			var input v1.InputMessage
			if err := conn.ReadJSON(&input); err != nil {
				fmt.Println(err)
				conn.Close()

				break
			}

			msg, err := input.ToEntity(r.Context())
			if err != nil {
				fmt.Println(err)
				conn.Close()

				break
			}

			if err := h.u.CreateMessage(context.Background(), &msg); err != nil {
				fmt.Println(err)
				conn.Close()

				break
			}
		}
	}()
}
