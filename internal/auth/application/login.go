package application

import (
	"context"

	domain "gcp-serverless-app/internal/auth/domain"
)

// LoginInput define los datos de entrada para el login.
type LoginInput struct {
	Email    string
	Password string
}

// LoginUseCase maneja la autenticación de usuarios.
type LoginUseCase struct {
	repo         domain.AuthRepository
	userProvider domain.UserProvider
	hasher       domain.PasswordHasher
}

// NewLoginUseCase crea una nueva instancia del caso de uso.
func NewLoginUseCase(repo domain.AuthRepository, userProvider domain.UserProvider, hasher domain.PasswordHasher) *LoginUseCase {
	return &LoginUseCase{repo: repo, userProvider: userProvider, hasher: hasher}
}

// Execute ejecuta el login del usuario.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*domain.UserInfo, error) {
	authUser, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	passwordHash := uc.hasher.Hash(input.Password, authUser.Salt)
	if passwordHash != authUser.PasswordHash {
		return nil, domain.ErrInvalidCredentials
	}

	userInfo, err := uc.userProvider.GetUser(ctx, authUser.UserID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
