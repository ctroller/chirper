package inject

import "github.com/jackc/pgx/v5/pgxpool"

type Application struct {
	DBPool *pgxpool.Pool
}

var App Application
