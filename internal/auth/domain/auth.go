package domain

import "time"

// AuthUser representa la entidad de autenticación del dominio.
type AuthUser struct {
	UserID       string
	Email        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// UserInfo es el DTO que retorna el login con datos del usuario.
type UserInfo struct {
	ID    string
	Email string
	Name  string
}
