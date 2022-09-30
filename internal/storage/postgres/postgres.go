package postgres

import (
	"github.com/gopher-market/internal/config"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
