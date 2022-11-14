package service

import (
	"github.com/gopher-market/internal/storage"
)

type BalanceService struct {
	repo storage.Orders
}

func NewBalanceService(repo storage.Orders) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) LoadOrder(orderID string) error {
	return nil
}

func (s *BalanceService) GetOrders(userID string) ([]string, error) {
	return nil, nil
}
