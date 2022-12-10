package handler

import (
	"context"
	"local/challengestockschat/stockschat/entity"
	"net/http"
	"time"

	v1 "local/challengestockschat/stockschat/gateway/http/handler/v1"

	"github.com/gorilla/websocket"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
	loggerPkg "github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var timeout = time.Minute

func (h Handler) Listen(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = erring.Wrap(err).Describe("failed to upgrade connection for listen request")

		httpresp.Internal(err).Handle(w, r)

		return
	}

	defer conn.Close()

	logger := loggerPkg.FromContext(r.Context())

	baseCtx, err := h.authorize(conn)
	if err != nil {
		_ = erring.Wrap(err).Describe("failed to authorize").Log(logger, zap.WarnLevel)

		return
	}

	baseCtx = loggerPkg.AddToCtx(baseCtx, logger)

	cleanup, err := h.u.Listen(baseCtx, func(m entity.Message) error {
		err := conn.WriteJSON(v1.ToMessangeOutput(m))
		if err != nil {
			return erring.Wrap(err).Log(logger, zap.ErrorLevel)
		}

		return nil
	})
	if err != nil {
		cleanup()

		_ = erring.Wrap(err).Log(logger, zap.ErrorLevel)

		return
	}

	defer cleanup()

	h.readMessages(baseCtx, conn)
}

func (h Handler) readMessages(ctx context.Context, conn *websocket.Conn) {
	logger := loggerPkg.FromContext(ctx)

	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()

		for {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			var input v1.InputMessage
			if err := conn.ReadJSON(&input); err != nil {
				_ = erring.Wrap(err).
					Describe("failed to read message").
					Log(logger, zap.ErrorLevel)

				conn.WriteJSON(v1.ToErrorOutput("invalid json"))

				break
			}

			if input.Content == "" {
				conn.WriteJSON(v1.ToErrorOutput("empty message not allowed"))
			}

			msg, err := input.ToEntity(ctx)
			if err != nil {
				_ = erring.Wrap(err).
					Describe("failed to encode message as json").
					Log(logger, zap.ErrorLevel)

				conn.WriteJSON(v1.ToErrorOutput("internal server error"))

				break
			}

			if err := h.u.CreateMessage(ctx, &msg); err != nil {
				_ = erring.Wrap(err).Log(logger, zap.ErrorLevel)

				conn.WriteJSON("failed to create message")

				break
			}
		}
	}()

	<-done
}

func (h Handler) authorize(conn *websocket.Conn) (context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var input v1.Authorization
	if err := conn.ReadJSON(&input); err != nil {
		conn.WriteJSON(v1.ToErrorOutput("invalid json"))

		return nil, erring.Wrap(err).
			Describe("failed to read message")
	}

	claims, err := h.authorizer.IsAuthorized(ctx, input.Token)
	if err != nil {
		conn.WriteJSON(v1.ToErrorOutput("authorization failed"))

		return nil, err
	}

	return handler.AddClaimsToContext(context.Background(), claims), nil
}
