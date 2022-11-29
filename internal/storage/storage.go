package storage

import (
	"context"
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
	SaveAccrual(order models.AccrualOrder) error
	GetOrders(userID int) (*[]models.Order, error)
	SaveOrderBalance(ctx context.Context, userID string, current float32) error
}

type Balance interface {
	GetBalance(ctx context.Context, userID string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawals models.Withdrawals) error // запрос на списание
	GetWithdrawals(userID string) ([]models.Withdrawals, error)         // информация о списаниях
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
