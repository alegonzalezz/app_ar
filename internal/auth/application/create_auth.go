package application

import (
	"context"
	"fmt"
	"time"

	domain "gcp-serverless-app/internal/auth/domain"
)

// CreateAuthInput define los datos de entrada para crear credenciales.
type CreateAuthInput struct {
	UserID    string
	Email     string
	Password  string
	CreatedAt time.Time
}

// CreateAuthUseCase maneja la creación de credenciales de autenticación.
type CreateAuthUseCase struct {
	repo   domain.AuthRepository
	hasher domain.PasswordHasher
}

// NewCreateAuthUseCase crea una nueva instancia del caso de uso.
func NewCreateAuthUseCase(repo domain.AuthRepository, hasher domain.PasswordHasher) *CreateAuthUseCase {
	return &CreateAuthUseCase{repo: repo, hasher: hasher}
}

// Execute ejecuta la creación de credenciales.
func (uc *CreateAuthUseCase) Execute(ctx context.Context, input CreateAuthInput) error {
	if input.Password == "" {
		return fmt.Errorf("la contraseña no puede estar vacía")
	}

	// Generar el salt a partir de la fecha de creación en formato RFC3339Nano
	salt := input.CreatedAt.Format(time.RFC3339Nano)

	// Hashear la contraseña concatenada con el salt
	passwordHash := uc.hasher.Hash(input.Password, salt)

	authUser := &domain.AuthUser{
		UserID:       input.UserID,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.CreatedAt,
	}

	return uc.repo.Save(ctx, authUser)
}
