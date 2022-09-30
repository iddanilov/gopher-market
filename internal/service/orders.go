package service

import (
	"github.com/gopher-market/internal/storage"
)

type OrderService struct {
	repo storage.Orders
}

func NewOrderService(repo storage.Orders) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) SaveOrder(orderID string) error {
	return nil
}

func (s *OrderService) GetOrders(userID string) ([]string, error) {
	return nil, nil
}
