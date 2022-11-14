package models

import "time"

type Order struct {
	Id         int        `json:"-" db:"user_id" storage:"id"`
	Number     string     `json:"number" db:"order_number" binding:"required"`
	Status     string     `json:"status" db:"status" binding:"required"`
	Accrual    string     `json:"accrual,omitempty" db:"accrual" `
	UploadedAt *time.Time `json:"uploaded_at" db:"uploaded_at"`
}
