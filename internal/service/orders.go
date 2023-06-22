package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/pkg/logging"
)

type OrderService struct {
	ctx       context.Context
	repo      storage.Orders
	logger    *logging.Logger
	loyaltyCh chan models.LoyaltyChan
}

func NewOrderService(ctx context.Context, repo storage.Orders, logger *logging.Logger, loyaltyCh chan models.LoyaltyChan) *OrderService {
	return &OrderService{ctx: ctx, repo: repo, logger: logger, loyaltyCh: loyaltyCh}
}

func (s *OrderService) LoadOrder(userID int, orderID string) (int, error) {
	order, err := s.repo.GetOrderByUserID(s.ctx, orderID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return http.StatusInternalServerError, err
		}
	}
	if order != nil {
		if order.ID != userID {
			return http.StatusConflict, errors.New("order already loaded")
		}
		return http.StatusOK, nil
	}

	err = s.repo.LoadOrder(userID, orderID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	s.loyaltyCh <- models.LoyaltyChan{
		UserID:  userID,
		OrderID: orderID,
	}
	return http.StatusAccepted, nil
}

func (s *OrderService) GetOrders(userID int) (*[]models.Order, error) {
	return s.repo.GetOrders(userID)
}
