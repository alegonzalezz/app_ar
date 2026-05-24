package auth

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
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
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
}
