package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/encryptor"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/repository"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/httpresp"
)

type Handler struct {
	handler.Handler
}

func NewHandler(p *pgxpool.Pool, jwtSecret string) Handler {
	encry := encryptor.NewEncryptor(jwtSecret)
	repo := repository.NewUserRepository(p)
	u := usecase.NewUsecase(encry, repo)

	return Handler{Handler: handler.New(u)}
}

func (h Handler) SetupRoutes(r chi.Router) {
	r.Post("/signup", httpresp.Handle(h.PostUser))
	r.Post("/login", httpresp.Handle(h.Login))
}
