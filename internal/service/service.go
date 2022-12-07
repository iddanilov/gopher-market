package service

import (
	"context"
	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/pkg/logging"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Orders interface {
	LoadOrder(userID int, orderID string) (int, error)
	GetOrders(userID int) (*[]models.Order, error)
}

type Balance interface {
	Withdraw(withdrawals models.Withdrawals) error           // запрос на списание
	GetWithdrawals(userID int) ([]models.Withdrawals, error) // информация о списаниях
	GetBalance(userID string) (models.Balance, error)
}

type Service struct {
	Authorization
	Orders
	Balance
	cfg    *config.Config
	logger *logging.Logger
}

func NewService(ctx context.Context, repos *storage.Storage, cfg *config.Config, logger *logging.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(ctx, repos.Authorization, logger),
		Orders:        NewOrderService(ctx, repos.Orders, *cfg, logger),
		Balance:       NewBalanceService(ctx, repos.Balance, logger),
	}
}
