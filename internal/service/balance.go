package service

import (
	"context"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/pkg/logging"
	"strconv"
)

type BalanceService struct {
	ctx    context.Context
	repo   storage.Balance
	logger *logging.Logger
}

func NewBalanceService(ctx context.Context, repo storage.Balance, logger *logging.Logger) *BalanceService {
	return &BalanceService{ctx: ctx, repo: repo, logger: logger}
}

func (s *BalanceService) GetBalance(userID string) (models.Balance, error) {
	return s.repo.GetBalance(userID)
}

func (s *BalanceService) Withdraw(withdrawals models.Withdrawals) error {
	return s.repo.Withdraw(withdrawals)
}
func (s *BalanceService) GetWithdrawals(userID int) ([]models.Withdrawals, error) {
	return s.repo.GetWithdrawals(strconv.Itoa(userID))
}
