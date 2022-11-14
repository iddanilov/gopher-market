package postgres

import (
	"github.com/jmoiron/sqlx"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

type Balance interface {
	LoadOrder(orderID string) error
	GetOrders(username string) ([]string, error)
}

func (r *BalancePostgres) LoadOrder(orderID string) error {
	return nil
}

func (r *BalancePostgres) GetOrders(username string) ([]string, error) {
	return []string{}, nil
}
