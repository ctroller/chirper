package db

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timtoronto634/pgx-slog"
)

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
	onceError  error
)

func NewPG(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			onceError = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		config.ConnConfig.Tracer = pgxslog.NewTracer(slog.Default())

		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			onceError = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		onceError = db.Ping(ctx)
		pgInstance = &Postgres{db}

	})

	return pgInstance, onceError
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
