package users

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmailAlreadyExists = errors.New("el correo electrónico ya se encuentra registrado")
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *User) error
}

type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuthCreator interface {
	CreateAuth(ctx context.Context, userID, email, password string, createdAt time.Time) error
}
