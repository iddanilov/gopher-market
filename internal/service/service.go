package service

import (
	"context"
	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Orders interface {
	LoadOrder(userID int, orderID string) error
	GetOrders(userID int) ([]models.Order, error)
}

type Balance interface {
	Withdraw(ctx context.Context, withdrawals models.Withdrawals) error // запрос на списание
	GetWithdrawals(userID int) ([]models.Withdrawals, error)            // информация о списаниях
	GetBalance(ctx context.Context, userID string) (models.Balance, error)
}

type Service struct {
	Authorization
	Orders
	Balance
	cfg *config.Config
}

func NewService(repos *storage.Storage, cfg *config.Config) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Orders:        NewOrderService(repos.Orders, *cfg),
		Balance:       NewBalanceService(repos.Balance),
	}
}
