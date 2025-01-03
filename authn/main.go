package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctroller/chirper/authn/internal/login"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func main() {
	setupLogging()
	r := setupRouter()

	slog.Info("Authentication service is running on port 5000")
	err := http.ListenAndServe(":5000", r)
	if err != nil {
		slog.Error("Failed to start authentication service", err)
	}
}
