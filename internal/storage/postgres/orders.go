package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gopher-market/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (r *OrdersPostgres) SaveAccrual(order models.AccrualOrder) error {
	query := `
UPDATE orders 
SET status = $1, accrual = $2 
WHERE order_number = $3`
	_, err := r.db.Exec(query, order.Status, order.Accrual, order.Order)
	if err != nil {
		log.Println("Can't Save accrual: ", err)
	}
	return err
}

func (r *OrdersPostgres) SaveOrderBalance(ctx context.Context, userID string, current int) error {
	var balance models.Balance

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {

		return fmt.Errorf("can't start tx; error: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				log.Println(ctx, err.Error())
			}
		}
	}(tx)

	if err := tx.GetContext(ctx, &balance, `
SELECT  user_id, user_current, withdrawn
FROM balance
WHERE user_id = $1 FOR UPDATE SKIP LOCKED
LIMIT 1`,
		userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(ctx, err.Error())
		} else {
			log.Println(ctx, err.Error())
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO balance (user_id, user_current, withdrawn)
VALUES($1, $2, $3)
    ON CONFLICT (user_id)
DO
UPDATE SET user_current=$2;`, userID, current+balance.Current, balance.Withdrawn)
	if err != nil {
		log.Println(ctx, err.Error())
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println(ctx, err.Error())
		return err
	}

	return nil
}

func (r *OrdersPostgres) LoadOrder(userID int, orderID string) error {
	query := `
INSERT INTO orders(order_number, user_id, status, accrual) 
VALUES ($1, $2, $3, $4)`
	log.Println(query)

	_, err := r.db.Exec(query, orderID, userID, "NEW", "0")
	if err != nil {
		log.Println("Can't Update Order: ", err)
	}
	return nil
}

func (r *OrdersPostgres) GetOrders(userID int) ([]models.Order, error) {
	var result []models.Order
	fmt.Println(userID)
	query := `
SELECT order_number, user_id, status, accrual, uploaded_at 
FROM orders
WHERE user_id = $1`
	if err := r.db.Select(&result, query, strconv.Itoa(userID)); err != nil {
		return nil, err
	}
	fmt.Println(result)
	return result, nil

}