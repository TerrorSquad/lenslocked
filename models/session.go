package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/terrorsquad/lenslocked/rand"
)

const (
	// MinBytesPerToken The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
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

func (sessionService *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := sessionService.BytesPerToken
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
		TokenHash: sessionService.hash(token),
	}

	row := sessionService.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2) ON CONFLICT (user_id) DO
			UPDATE SET token_hash = $2 RETURNING id;`,
		session.UserID, session.TokenHash,
	)
	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (sessionService *SessionService) User(token string) (*User, error) {
	tokenHash := sessionService.hash(token)
	row := sessionService.DB.QueryRow(
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

func (sessionService *SessionService) Delete(token string) error {
	tokenHash := sessionService.hash(token)
	_, err := sessionService.DB.Exec(
		`DELETE FROM sessions WHERE token_hash = $1;`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (sessionService *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
