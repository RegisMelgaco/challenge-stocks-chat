package handler

import (
	"local/challengestockschat/stockschat/gateway/postgres/repository"
	"local/challengestockschat/stockschat/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct {
	u        usecase.Usecase
	upgrader websocket.Upgrader
}

func NewHandler(pool *pgxpool.Pool, broker usecase.Broker) Handler {
	repo := repository.New(pool)
	u := usecase.New(repo, broker)

	return Handler{u, websocket.Upgrader{}}
}

func (h Handler) Route(r chi.Router) {
	r.Get("/listen", h.Listen)
}
