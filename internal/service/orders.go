package service

import (
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
)

type OrderService struct {
	repo storage.Orders
}

func NewOrderService(repo storage.Orders) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) LoadOrder(userID int, orderID string) error {
	return s.repo.LoadOrder(userID, orderID)
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	result, err := s.repo.GetOrders(userID)
	return result, err
}
