package postgres

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/b0pof/avito-internship/internal/config"
)

func NewPgxDatabase(ctx context.Context, cfg config.Postgres) (*sqlx.DB, error) {
	dbClient, err := sqlx.ConnectContext(ctx, "pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}
