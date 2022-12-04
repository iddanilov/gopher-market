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
	"io"
	"log"
	"net/http"
	"strconv"
)

type OrderService struct {
	repo storage.Orders
	cfg  config.Config
}

func NewOrderService(repo storage.Orders, cfg config.Config) *OrderService {
	return &OrderService{repo: repo, cfg: cfg}
}

func (s *OrderService) LoadOrder(userID int, orderID string) error {
	ctx := context.Background()
	order, err := s.repo.GetOrderByUserID(ctx, userID, orderID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	if order != nil {
		return errors.New("order already loaded")
	}

	err = s.repo.LoadOrder(userID, orderID)
	if err != nil {
		return err
	}
	go s.writeInfoAboutOrder(userID, orderID)
	return nil
}

func (s *OrderService) GetOrders(userID int) (*[]models.Order, error) {
	return s.repo.GetOrders(userID)
}

func (s *OrderService) writeInfoAboutOrder(userID int, orderID string) {
	ctx := context.Background()
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
		log.Println("ERROR can't save accrual: ", err)
	}

	go func() {
		err = s.repo.SaveAccrual(order)
		if err != nil {
			log.Println("ERROR can't save accrual: ", err)
		}
	}()
	if err != nil {
		log.Println("ERROR can't save accrual: ", err)
	}
	go func() {
		err = s.repo.SaveOrderBalance(ctx, strconv.Itoa(userID), order.Accrual)
		if err != nil {
			log.Println("ERROR can't save accrual: ", err)
		}
	}()

}
