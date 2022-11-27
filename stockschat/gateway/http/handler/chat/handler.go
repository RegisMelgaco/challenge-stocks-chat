package chat

import (
	chatRepository "local/challengestockschat/stockschat/gateway/postgres/repository/chat"
	"local/challengestockschat/stockschat/usecase/chat"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/httpresp"
)

type Handler struct {
	u chat.Usecase
}

func NewHandler(pool *pgxpool.Pool) Handler {
	repo := chatRepository.NewRepository(pool)
	u := chat.NewUsecase(repo)

	return Handler{u}
}

func (h Handler) Route(r chi.Router) {
	r.Post("/message", httpresp.Handle(h.PostMessage))
}
