package models

type Balance struct {
	Id        int    `json:"-" storage:"id"`
	Current   string `json:"current" binding:"required"`
	Withdrawn string `json:"withdrawn" binding:"required"`
}
