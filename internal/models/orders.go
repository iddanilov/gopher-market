package models

import "time"

type Order struct {
	ID         int        `json:"-" db:"user_id" storage:"id"`
	Number     string     `json:"number" db:"order_number" binding:"required"`
	Status     string     `json:"status" db:"status" binding:"required"`
	Accrual    float32    `json:"accrual,omitempty" db:"accrual" `
	UploadedAt *time.Time `json:"uploaded_at" db:"uploaded_at"`
}

type AccrualOrder struct {
	Order      string     `json:"order" db:"order_number" binding:"required"`
	Status     string     `json:"status" db:"status" binding:"required"`
	Accrual    float32    `json:"accrual,omitempty" db:"accrual" `
	UploadedAt *time.Time `json:"uploaded_at" db:"uploaded_at"`
}
