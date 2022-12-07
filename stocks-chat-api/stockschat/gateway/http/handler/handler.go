package handler

import (
	"local/challengestockschat/stockschat/usecase"
	chatUsecase "local/challengestockschat/stockschat/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	authUsecase "github.com/regismelgaco/go-sdks/auth/auth/usecase"
)

type Handler struct {
	u          chatUsecase.Usecase
	upgrader   websocket.Upgrader
	authorizer authUsecase.Usecase
}

func NewHandler(u usecase.Usecase, authorizer authUsecase.Usecase) Handler {
	return Handler{u, websocket.Upgrader{}, authorizer}
}

func (h Handler) Route(r chi.Router) {
	r.Get("/listen", h.Listen)
}
