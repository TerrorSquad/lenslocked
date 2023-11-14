package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/terrorsquad/lenslocked/rand"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     uint
	UserID uint
	// Token is only set when a password reset is created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes should be used when generating
	// each password reset token. If this value is not set or is less than the MinBytesPerToken const it will be ignored
	// and the MinBytesPerToken value will be used instead.
	BytesPerToken int
	// Duration is the amount of time that a password reset is valid for.
	// Default to DefaultResetDuration.
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userId uint
	row := service.DB.QueryRow(`SELECT id FROM users WHERE email = $1;`, email)
	err := row.Scan(&userId)
	if err != nil {
		// TODO: Consider returning a specific error if the user does not exist.
		return nil, fmt.Errorf("create password reset: %w", err)
	}
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	passwordReset := PasswordReset{
		UserID:    userId,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = service.DB.QueryRow(`
			INSERT INTO password_resets (user_id, token_hash, expires_at)
			VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
			UPDATE SET token_hash = $2 RETURNING id;`,
		passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt,
	)
	err = row.Scan(&passwordReset.ID)

	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	return &passwordReset, nil

}

func (service *PasswordResetService) Consume(token string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume()")
}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
