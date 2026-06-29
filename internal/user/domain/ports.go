package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmailAlreadyExists = errors.New("el correo electrónico ya se encuentra registrado")
)

// UserRepository define el puerto driven para acceso a datos de usuario.
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) error
}

// TxManager define el puerto driven para gestión de transacciones.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// AuthCreator define el puerto driven para creación de credenciales.
// Implementado por un bridge inter-módulo.
type AuthCreator interface {
	CreateAuth(ctx context.Context, authID, profileID, profileType, email, password string, createdAt time.Time) error
}
