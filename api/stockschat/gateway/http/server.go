package http

import (
	"local/challengestockschat/stockschat/config"
	"local/challengestockschat/stockschat/gateway/http/handler"
	"local/challengestockschat/stockschat/usecase"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/encryptor"
	authHandler "github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/repository"
	authUsecase "github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/httpresp"
	"go.uber.org/zap"
)

const timeout = 5 * time.Minute

func SetupRouter(logger *zap.Logger, pool *pgxpool.Pool, cfg config.Config, usecase usecase.Usecase) chi.Router {
	r := chi.NewRouter()

	authorizer := authUsecase.NewUsecase(
		encryptor.NewEncryptor(cfg.JWTSecret),
		repository.NewUserRepository(pool),
	)
	chatHandler := handler.NewHandler(usecase, authorizer, cfg.IsDev)
	authHandler := authHandler.NewHandler(authorizer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"Authorization",
		},
		AllowCredentials: true,
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader(
		"Content-type", "application/json",
	))
	r.Use(httpresp.Logger(logger))
	r.Use(httpresp.RequestID)

	r.With(middleware.Timeout(timeout)).Route("/api/auth", authHandler.SetupRoutes)
	r.Route("/api/chat", chatHandler.Route)

	r.Get("/api/healthcheck", func(w http.ResponseWriter, r *http.Request) {})

	return r
}
