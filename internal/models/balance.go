package models

import "time"

type Balance struct {
	ID        string `json:"-" db:"user_id"`
	Current   int    `json:"current" db:"user_current" binding:"required"`
	Withdrawn int    `json:"withdrawn" db:"withdrawn" binding:"required"`
}

type Withdrawals struct {
	ID          string     `json:"-" db:"user_id"`
	Order       string     `json:"order" db:"order_number"`
	Sum         int        `json:"sum" db:"sum"`
	ProcessedAt *time.Time `json:"processed_at" db:"processed_at"`
}
