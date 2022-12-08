package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gopher-market/internal/config"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/internal/storage"
	"github.com/gopher-market/pkg/logging"
)

type LoyaltyService struct {
	ctx       context.Context
	repo      storage.Orders
	cfg       config.Config
	logger    *logging.Logger
	loyaltyCh chan models.LoyaltyChan
	client    *http.Client
}

func NewLoyaltyService(ctx context.Context, repo storage.Orders, cfg config.Config, logger *logging.Logger, loyaltyCh chan models.LoyaltyChan) *LoyaltyService {
	return &LoyaltyService{ctx: ctx, repo: repo, cfg: cfg, logger: logger, loyaltyCh: loyaltyCh, client: &http.Client{}}
}

func (s *LoyaltyService) GetLoyalty() {
	for loyalty := range s.loyaltyCh {
		go func(loyalty models.LoyaltyChan) {
			order := models.AccrualOrder{}
			err := s.sendLoyaltyRequest(&order, loyalty)
			if err != nil {
				log.Println("ERROR can't save accrual: ", err)
				return
			}
			go func(order models.AccrualOrder) {
				err = s.repo.SaveAccrual(order)
				if err != nil {
					log.Println("ERROR can't save accrual: ", err)
					return
				}
			}(order)
			if err != nil {
				log.Println("ERROR can't save accrual: ", err)
			}
			go func(userID string, accrual float32) {
				err = s.repo.SaveOrderBalance(userID, accrual)
				if err != nil {
					log.Println("ERROR can't save accrual: ", err)
				}
			}(strconv.Itoa(loyalty.UserID), order.Accrual)
		}(loyalty)
	}
}

func (s *LoyaltyService) sendLoyaltyRequest(order *models.AccrualOrder, loyalty models.LoyaltyChan) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", s.cfg.Accrual.Address, loyalty.OrderID), nil)
	res, err := s.client.Do(req)
	if res.StatusCode != http.StatusOK {
		return err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	err = json.Unmarshal(body, &order)
	if err != nil {
		s.logger.Error("ERROR can't save accrual: ", err)
		return err
	}
	return nil
}
