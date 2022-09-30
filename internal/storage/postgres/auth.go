package postgres

import (
	"fmt"
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
	query := fmt.Sprintf("INSERT INTO users(login, password_hash) values ($1, $2)")
	log.Println(query)
	log.Println(models)

	_, err := r.db.Exec(query, models.Login, models.Password)
	if err != nil {
		log.Println("Can't Update Metric: ", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := "SELECT id FROM users WHERE login=$1 AND password_hash=$2"
	err := r.db.Get(&user, query, login, password)

	return user, err
}
