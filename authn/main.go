package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctroller/chirper/authn/internal/inject"
	"github.com/ctroller/chirper/authn/internal/login"
	"github.com/ctroller/chirper/authn/lib/db"
	"github.com/ctroller/chirper/authn/lib/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ctroller/chirper/authn/docs"
)

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

func setupDB() *db.Postgres {
	postgres, err := db.NewPG(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to create DB connection", "error", err)
		os.Exit(1)
	}

	var dbName string
	err = postgres.DB.QueryRow(context.Background(), "SELECT current_database();").Scan(&dbName)
	if err != nil {
		slog.Error("Unable to connect to DB", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Connected to DB", "database", dbName)
	}

	return postgres
}

func setupApp() {
	db := setupDB()

	inject.App = inject.Application{
		DBPool: db.DB,
		UserRepository: &user.UserRepositoryImpl{
			DB: db.DB,
		},
	}
}

func main() {
	setupLogging()
	setupApp()
	r := setupRouter()

	slog.Info("Authentication service is running on port 5000")
	err := http.ListenAndServe(":5000", r)
	if err != nil {
		slog.Error("Failed to start authentication service", "error", err)
	}
}
