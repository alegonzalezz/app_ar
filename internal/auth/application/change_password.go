package application

import (
	"context"

	domain "gcp-serverless-app/internal/auth/domain"
)

// ChangePasswordInput define los datos de entrada para cambiar contraseña.
type ChangePasswordInput struct {
	Email       string
	OldPassword string
	NewPassword string
}

// ChangePasswordUseCase maneja el cambio de contraseña.
type ChangePasswordUseCase struct {
	repo   domain.AuthRepository
	hasher domain.PasswordHasher
}

// NewChangePasswordUseCase crea una nueva instancia del caso de uso.
func NewChangePasswordUseCase(repo domain.AuthRepository, hasher domain.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo, hasher: hasher}
}

// Execute ejecuta el cambio de contraseña.
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
	authUser, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	currentHash := uc.hasher.Hash(input.OldPassword, authUser.Salt)
	if currentHash != authUser.PasswordHash {
		return domain.ErrInvalidCredentials
	}

	newHash := uc.hasher.Hash(input.NewPassword, authUser.Salt)

	return uc.repo.UpdatePassword(ctx, authUser.UserID, newHash)
}
