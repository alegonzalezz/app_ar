package domain

import "context"

// AuthRepository define el puerto driven para persistencia de autenticación.
type AuthRepository interface {
	Save(ctx context.Context, authUser *AuthUser) error
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
}

// UserProvider define el puerto driven para obtener datos del usuario.
// Implementado por un bridge inter-módulo.
type UserProvider interface {
	GetUser(ctx context.Context, userID string) (*UserInfo, error)
}

// PasswordHasher define el puerto driven para hashing de contraseñas.
type PasswordHasher interface {
	Hash(password, salt string) string
}
