package storage

import (
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Orders interface {
	SaveOrder(orderID string) error
	GetOrders(username string) ([]string, error)
}

type Storage struct {
	Authorization
	Orders
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(db),
	}
}
