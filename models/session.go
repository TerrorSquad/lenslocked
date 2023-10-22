package models

type Session struct {
	ID        uint
	UserID    uint
	TokenHash string
}
