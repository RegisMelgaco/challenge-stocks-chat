package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	authUsecase "github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"local/challengestockschat/stockschat/usecase"
	chatUsecase "local/challengestockschat/stockschat/usecase"
)

type Handler struct {
	u          chatUsecase.Usecase
	upgrader   websocket.Upgrader
	authorizer authUsecase.Usecase
}

func NewHandler(u usecase.Usecase, authorizer authUsecase.Usecase, isDev bool) Handler {
	upgrader := websocket.Upgrader{}

	if isDev {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}

	return Handler{u, upgrader, authorizer}
}

func (h Handler) Route(r chi.Router) {
	r.Get("/listen", h.Listen)
}
