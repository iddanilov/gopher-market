package models

type User struct {
	ID       int    `json:"-" storage:"id"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
