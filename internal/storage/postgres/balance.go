package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gopher-market/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r *BalancePostgres) Withdraw(ctx context.Context, withdrawals models.Withdrawals) error {
	var balance models.Balance

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {

		return errors.New(fmt.Sprintf("can't start tx; error: %v", err))
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
		withdrawals.Id); err != nil {
		log.Println(ctx, err.Error())
		return err

	}

	if balance.Current < withdrawals.Sum {
		return errors.New(fmt.Sprintf("insufficient funds! on ammount: %v", balance.Current))
	}

	_, err = tx.ExecContext(ctx, `
UPDATE balance 
SET user_current = $2, withdrawn = $3 
WHERE user_id = $1`,
		withdrawals.Id, balance.Current-withdrawals.Sum, balance.Withdrawn+withdrawals.Sum,
	)
	if err != nil {
		{
			err := tx.Rollback()
			if err != sql.ErrTxDone {
				return err
			}
		}

		log.Println(ctx, err.Error())
		return err
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO withdrawals(user_id, order_number, sum) 
VALUES ($1, $2, $3)`,
		withdrawals.Id, withdrawals.Order, withdrawals.Sum,
	)
	if err != nil {
		{
			err := tx.Rollback()
			if err != sql.ErrTxDone {
				return err
			}
		}

		log.Println(ctx, err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(ctx, err.Error())
		{
			err := tx.Rollback()
			if err != sql.ErrTxDone {
				return err
			}
		}
		return err
	}

	return nil
}

func (r *BalancePostgres) GetWithdrawals(userID string) ([]models.Withdrawals, error) {
	var balance []models.Withdrawals
	query := fmt.Sprintf(`
		SELECT user_id, order_number, sum, processed_at 
		FROM withdrawals
    	WHERE user_id = $1`)
	if err := r.db.Select(&balance, query, userID); err != nil {
		return nil, err
	}

	return balance, nil
}

func (r *BalancePostgres) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
	balance := models.Balance{}

	if err := r.db.GetContext(ctx, &balance, `
SELECT  user_id, user_current, withdrawn
FROM balance
WHERE user_id = $1 FOR UPDATE SKIP LOCKED
LIMIT 1`,
		userID); err != nil {

		log.Println(ctx, err.Error())
		return balance, err
	}

	return balance, nil
}
