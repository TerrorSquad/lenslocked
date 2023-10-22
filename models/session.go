package models

import "database/sql"

type Session struct {
	ID     uint
	UserID uint
	// Token is only set when creating a new session. When looking up a session
	// this will be empty, as we only store the hashed version of the token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID uint) (*Session, error) {
	// TODO: 1. Create the session token
	// TODO: 2. Implement SessionService.Create
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
