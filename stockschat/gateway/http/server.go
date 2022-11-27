package http

import (
	"encoding/json"
	"local/challengestockschat/stockschat/config"
	"local/challengestockschat/stockschat/gateway/http/handler/chat"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	auth "github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/httpresp"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger, pool *pgxpool.Pool, cfg config.Config) chi.Router {
	r := chi.NewRouter()

	authHandler := auth.NewHandler(pool, cfg.JWTSecret)
	chatHandler := chat.NewHandler(pool)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader(
		"Content-type", "application/json",
	))
	r.Use(httpresp.Logger(logger))

	r.Route("/auth", authHandler.SetupRoutes)
	r.With(authHandler.IsAuthorized).Route("/chat", chatHandler.Route)

	r.With(authHandler.IsAuthorized).Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"msg": "pong"})
	})

	return r
}
