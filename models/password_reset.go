package models

import (
	"database/sql"
	"fmt"
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
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Create()")
}

func (service *PasswordResetService) Consume(token string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: implement PasswordResetService.Consume()")
}
