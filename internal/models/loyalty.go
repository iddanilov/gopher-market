package models

import "time"

type LoyaltyChan struct {
	UserID  int
	OrderID string
}

type Loyalty struct {
	Order      string     `json:"order" db:"order_number" binding:"required"`
	Status     string     `json:"status" db:"status" binding:"required"`
	Accrual    float32    `json:"accrual,omitempty" db:"accrual" `
	UploadedAt *time.Time `json:"uploaded_at" db:"uploaded_at"`
}
