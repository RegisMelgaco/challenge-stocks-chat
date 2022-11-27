package http

import (
	"encoding/json"
	"local/challengestockschat/stockschat/config"
	"local/challengestockschat/stockschat/gateway/http/handler/auth"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, cfg config.Config) chi.Router {
	r := chi.NewRouter()

	authHandler := auth.NewHandler(pool, cfg.JWTSecret)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader(
		"Content-type", "application/json",
	))

	r.Route("/auth", authHandler.SetupRoutes)

	r.With(authHandler.IsAuthorized).Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"msg": "pong"})
	})

	return r
}
