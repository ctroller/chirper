package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctroller/chirper/authn/internal/inject"
	"github.com/ctroller/chirper/authn/internal/login"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ctroller/chirper/authn/docs"
)

// gin-swagger middleware
// swagger embed files

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/api/v1/login", login.LoginHandler)

	r.Mount("/swagger", httpSwagger.WrapHandler)

	return r
}

func setupLogging() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func setupDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	return dbpool
}

func setupApp() {
	inject.App = inject.Application{
		DBPool: setupDB(),
	}
}

func main() {
	setupApp()
	setupLogging()
	r := setupRouter()

	slog.Info("Authentication service is running on port 5000")
	err := http.ListenAndServe(":5000", r)
	if err != nil {
		slog.Error("Failed to start authentication service", err)
	}
}
