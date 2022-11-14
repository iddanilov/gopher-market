package service

import (
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
	LoadOrder(orderID string) error
	GetOrders(userID string) ([]string, error)
}

type Service struct {
	Authorization
	Orders
	Balance
}

func NewService(repos *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Orders:        NewOrderService(repos.Orders),
		//Balance:       NewBalanceService(repos.Balance),
	}
}
