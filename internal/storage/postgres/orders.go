package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/pkg/logging"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type OrdersPostgres struct {
	ctx    context.Context
	db     *sqlx.DB
	logger *logging.Logger
}

func NewOrdersPostgres(ctx context.Context, db *sqlx.DB, logger *logging.Logger) *OrdersPostgres {
	return &OrdersPostgres{ctx: ctx, db: db, logger: logger}
}

func (r *OrdersPostgres) SaveAccrual(order models.AccrualOrder) error {
	query := `
UPDATE orders 
SET status = $1, accrual = $2 
WHERE order_number = $3`
	_, err := r.db.ExecContext(r.ctx, query, order.Status, order.Accrual, order.Order)
	if err != nil {
		r.logger.Error("Can't Save accrual: ", err)
	}
	return err
}

func (r *OrdersPostgres) SaveOrderBalance(userID string, current float32) error {
	var balance models.Balance

	tx, err := r.db.BeginTxx(r.ctx, nil)
	if err != nil {
		return fmt.Errorf("can't start tx; error: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				r.logger.Debug(r.ctx, err.Error())
			}
		}
	}(tx)

	if err := tx.GetContext(r.ctx, &balance, `
SELECT  user_id, user_current, withdrawn
FROM balance
WHERE user_id = $1 FOR UPDATE SKIP LOCKED
LIMIT 1`,
		userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug(r.ctx, err.Error())
		} else {
			r.logger.Error(r.ctx, err.Error())
			return err
		}
	}

	_, err = tx.ExecContext(r.ctx, `
INSERT INTO balance (user_id, user_current, withdrawn)
VALUES($1, $2, $3)
    ON CONFLICT (user_id)
DO
UPDATE SET user_current=$2;`, userID, current+balance.Current, balance.Withdrawn)
	if err != nil {
		r.logger.Error(r.ctx, err.Error())
		return err
	}
	if err := tx.Commit(); err != nil {
		r.logger.Error(r.ctx, err.Error())
		return err
	}

	return nil
}

func (r *OrdersPostgres) LoadOrder(userID int, orderID string) error {

	query := `
INSERT INTO orders(order_number, user_id, status, accrual) 
VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(r.ctx, query, orderID, userID, "NEW", "0")
	if err != nil {
		r.logger.Error("Can't Update Order: ", err)
		return err
	}
	return nil
}

func (r *OrdersPostgres) GetOrderByUserID(ctx context.Context, orderID string) (*models.Order, error) {
	var result models.Order
	query := `
SELECT order_number, user_id, status, accrual, uploaded_at
FROM orders
WHERE order_number =$1
LIMIT 1`
	if err := r.db.GetContext(ctx, &result, query, orderID); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *OrdersPostgres) GetOrders(userID int) (*[]models.Order, error) {
	var result []models.Order
	query := `
SELECT order_number, user_id, status, accrual, uploaded_at 
FROM orders
WHERE user_id = $1`
	if err := r.db.SelectContext(r.ctx, &result, query, strconv.Itoa(userID)); err != nil {
		return nil, err
	}
	return &result, nil
}
