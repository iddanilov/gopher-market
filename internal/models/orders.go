package models

import "time"

type Order struct {
	Id         int           `json:"-" storage:"id"`
	Number     string        `json:"number" binding:"required"`
	Status     string        `json:"status" binding:"required"`
	Accrual    string        `json:"accrual"`
	UploadedAt time.Duration `json:"uploaded_at"`
}
