package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/pkg/logging"
	"io"
	"log"
	"net/http"
	"strconv"
)

type OrderService struct {
	ctx    context.Context
	repo   storage.Orders
	cfg    config.Config
	logger *logging.Logger
}

func NewOrderService(ctx context.Context, repo storage.Orders, cfg config.Config, logger *logging.Logger) *OrderService {
	return &OrderService{ctx: ctx, repo: repo, cfg: cfg, logger: logger}
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
	go s.writeInfoAboutOrder(userID, orderID)
	return http.StatusAccepted, nil
}

func (s *OrderService) GetOrders(userID int) (*[]models.Order, error) {
	return s.repo.GetOrders(userID)
}

func (s *OrderService) writeInfoAboutOrder(userID int, orderID string) {
	client := &http.Client{}

	order := models.AccrualOrder{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", s.cfg.Accrual.Address, orderID), nil)
	res, err := client.Do(req)

	if res.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	err = json.Unmarshal(body, &order)
	if err != nil {
		s.logger.Error("ERROR can't save accrual: ", err)
		return
	}

	go func() {
		err = s.repo.SaveAccrual(order)
		if err != nil {
			log.Println("ERROR can't save accrual: ", err)
			return
		}
	}()
	if err != nil {
		log.Println("ERROR can't save accrual: ", err)
	}
	go func() {
		err = s.repo.SaveOrderBalance(strconv.Itoa(userID), order.Accrual)
		if err != nil {
			log.Println("ERROR can't save accrual: ", err)
		}
	}()

}
