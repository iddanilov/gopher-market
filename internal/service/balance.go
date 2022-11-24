package service

import (
	"context"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"strconv"
)

type BalanceService struct {
	repo storage.Balance
}

func NewBalanceService(repo storage.Balance) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
	return s.repo.GetBalance(ctx, userID)
}

func (s *BalanceService) Withdraw(ctx context.Context, withdrawals models.Withdrawals) error {
	return s.repo.Withdraw(ctx, withdrawals)
}
func (s *BalanceService) GetWithdrawals(userID int) ([]models.Withdrawals, error) {
	return s.repo.GetWithdrawals(strconv.Itoa(userID))
}
