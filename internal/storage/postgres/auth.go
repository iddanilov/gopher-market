package postgres

import (
	"database/sql"
	"errors"
	"github.com/gopher-market/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(models models.User) (int, error) {
	var id int
	var passwordHash string
	query := "SELECT password_hash FROM users WHERE login=$1"
	err := r.db.Get(&passwordHash, query, models.Login)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("can't registered user: ", err)
			return 0, errors.New("can't registered user")
		}
	} else {
		if passwordHash != models.Password {
			log.Println("user already registered: ", err)
			return 0, errors.New("user already registered")
		}
	}

	query = `
INSERT INTO users(login, password_hash)
VALUES ($1, $2)`
	log.Println(query)
	log.Println(models)

	_, err = r.db.Exec(query, models.Login, models.Password)
	if err != nil {
		log.Println("can't registered user: ", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := "SELECT id FROM users WHERE login=$1 AND password_hash=$2"
	err := r.db.Get(&user, query, login, password)

	return user, err
}
