package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage/postgres"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Orders interface {
	LoadOrder(userID int, orderID string) error
	GetOrders(userID int) ([]models.Order, error)
}

type Balance interface {
	//LoadOrder(orderID string) error
	//GetOrders(username string) ([]string, error)
}

type Storage struct {
	Authorization
	Orders
	Balance
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(db),
		Orders:        postgres.NewOrdersPostgres(db),
		Balance:       postgres.NewBalancePostgres(db),
	}
}
