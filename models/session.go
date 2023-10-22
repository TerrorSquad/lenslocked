package models

import (
	"database/sql"
	"fmt"
	"github.com/terrorsquad/lenslocked/rand"
)

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
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	// TODO: Hash the session token
	session := Session{
		UserID: userID,
		Token:  token,
		// TODO: Set the token hash
	}
	// todo: Store the session in the DB

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
