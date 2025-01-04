package inject

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	DBPool *pgxpool.Pool
	Logger *slog.Logger
}

var App Application
