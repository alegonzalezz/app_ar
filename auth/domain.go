package auth

import (
	"context"
	"time"
)

type AuthUser struct {
	UserID       string     `json:"user_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Salt         string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type Repository interface {
	Save(ctx context.Context, authUser *AuthUser) error
}
