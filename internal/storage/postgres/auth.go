package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	ctx    context.Context
	db     *sqlx.DB
	logger *logging.Logger
}

func NewAuthPostgres(ctx context.Context, db *sqlx.DB, logger *logging.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, logger: logger, ctx: ctx}
}

func (r *AuthPostgres) CreateUser(models models.User) (int, error) {
	var id int
	var passwordHash string
	query := "SELECT password_hash FROM users WHERE login=$1"
	err := r.db.GetContext(r.ctx, &passwordHash, query, models.Login)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Debug("can't registered user: ", err)
			return 0, errors.New("can't registered user")
		}
	} else {
		if passwordHash != models.Password {
			r.logger.Error("user already registered: ", err)
			return 0, errors.New("user already registered")
		}
	}

	query = `
INSERT INTO users(login, password_hash)
VALUES ($1, $2)`

	_, err = r.db.ExecContext(r.ctx, query, models.Login, models.Password)
	if err != nil {
		r.logger.Error("can't registered user: ", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := "SELECT id FROM users WHERE login=$1 AND password_hash=$2"
	err := r.db.GetContext(r.ctx, &user, query, login, password)

	return user, err
}
