package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/terrorsquad/lenslocked/rand"
)

const (
	// MinBytesPerToken The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
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
	// BytesPerToken is used to determine how many bytes should be used when generating
	// each session token. If this value is not set or is less than the MinBytesPerToken const it will be ignored
	// and the MinBytesPerToken value will be used instead.
	BytesPerToken int
}

func (ss *SessionService) Create(userID uint) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// 1. try to update the session
	// 2. if no session exists, create a new one

	row := ss.DB.QueryRow(
		`UPDATE sessions SET token_hash = $2 WHERE user_id = $1 RETURNING id;`,
		session.UserID, session.TokenHash,
	)
	err = row.Scan(&session.ID)
	if errors.Is(err, sql.ErrNoRows) {
		row = ss.DB.QueryRow(
			`INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;`,
			session.UserID, session.TokenHash,
		)
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	row := ss.DB.QueryRow(
		`
				SELECT users.id, users.email, users.password_hash
				FROM users
					JOIN sessions ON users.id = sessions.user_id
				WHERE sessions.token_hash = $1;
			`,
		tokenHash,
	)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(
		`DELETE FROM sessions WHERE token_hash = $1;`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
