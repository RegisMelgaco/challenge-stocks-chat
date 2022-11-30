package chat

import (
	chatRepository "local/challengestockschat/stockschat/gateway/postgres/repository/chat"
	"local/challengestockschat/stockschat/usecase/chat"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct {
	u        chat.Usecase
	upgrader websocket.Upgrader
}

func NewHandler(pool *pgxpool.Pool, broker chat.Broker) Handler {
	repo := chatRepository.NewRepository(pool)
	u := chat.NewUsecase(repo, broker)

	return Handler{u, websocket.Upgrader{}}
}

func (h Handler) Route(r chi.Router) {
	r.Get("/listen", h.Listen)
}
