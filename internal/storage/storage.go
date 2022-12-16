package storage

import (
	"context"
	"github.com/gopher-market/pkg/logging"
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
	GetOrderByUserID(ctx context.Context, orderID string) (*models.Order, error)
	GetOrders(userID int) (*[]models.Order, error)
	SaveOrderBalance(userID string, current float32) error
}

type Balance interface {
	GetBalance(userID string) (models.Balance, error)
	Withdraw(withdrawals models.Withdrawals) error              // запрос на списание
	GetWithdrawals(userID string) ([]models.Withdrawals, error) // информация о списаниях
}

type Storage struct {
	logger *logging.Logger
	Authorization
	Orders
	Balance
}

func NewStorage(ctx context.Context, db *sqlx.DB, logger *logging.Logger) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(ctx, db, logger),
		Orders:        postgres.NewOrdersPostgres(ctx, db, logger),
		Balance:       postgres.NewBalancePostgres(ctx, db, logger),
	}
}
