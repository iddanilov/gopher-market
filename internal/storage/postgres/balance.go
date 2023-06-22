package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/gopher-market/internal/models"
	"github.com/gopher-market/pkg/logging"
)

type BalancePostgres struct {
	ctx    context.Context
	db     *sqlx.DB
	logger *logging.Logger
}

func NewBalancePostgres(ctx context.Context, db *sqlx.DB, logger *logging.Logger) *BalancePostgres {
	return &BalancePostgres{ctx: ctx, db: db, logger: logger}
}

func (r *BalancePostgres) Withdraw(withdrawals models.Withdrawals) error {
	var balance models.Balance

	tx, err := r.db.BeginTxx(r.ctx, nil)
	if err != nil {
		return fmt.Errorf("can't start tx; error: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil {
			r.logger.Error(r.ctx, err.Error())
		}
	}(tx)

	if err := tx.GetContext(r.ctx, &balance, `
SELECT  user_id, user_current, withdrawn
FROM balance
WHERE user_id = $1 FOR UPDATE SKIP LOCKED
LIMIT 1`,
		withdrawals.ID); err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(r.ctx, err.Error())
			return err
		}
	}

	if balance.Current < withdrawals.Sum {
		return fmt.Errorf("insufficient funds! on ammount: %v", balance.Current)
	}

	_, err = tx.ExecContext(r.ctx, `
UPDATE balance 
SET user_current = $2, withdrawn = $3 
WHERE user_id = $1`,
		withdrawals.ID, balance.Current-withdrawals.Sum, balance.Withdrawn+withdrawals.Sum,
	)
	if err != nil {
		{
			err := tx.Rollback()
			if err != sql.ErrTxDone {
				return err
			}
		}

		r.logger.Error(r.ctx, err.Error())
		return err
	}

	_, err = tx.ExecContext(r.ctx, `
INSERT INTO withdrawals(user_id, order_number, sum) 
VALUES ($1, $2, $3)`,
		withdrawals.ID, withdrawals.Order, withdrawals.Sum,
	)
	if err != nil {
		{
			err := tx.Rollback()
			if err != sql.ErrTxDone {
				return err
			}
		}

		r.logger.Error(r.ctx, err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error(r.ctx, err.Error())
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
	query := `
SELECT user_id, order_number, sum, processed_at 
FROM withdrawals
WHERE user_id = $1`
	if err := r.db.SelectContext(r.ctx, &balance, query, userID); err != nil {
		return nil, err
	}

	return balance, nil
}

func (r *BalancePostgres) GetBalance(userID string) (models.Balance, error) {
	balance := models.Balance{}

	if err := r.db.GetContext(r.ctx, &balance, `
SELECT  user_id, user_current, withdrawn
FROM balance
WHERE user_id = $1 FOR UPDATE SKIP LOCKED
LIMIT 1`,
		userID); err != nil {

		r.logger.Error(r.ctx, err.Error())
		return balance, err
	}

	return balance, nil
}
