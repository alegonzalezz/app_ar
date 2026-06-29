package application

import (
	"context"
	"fmt"
	"time"

	domain "gcp-serverless-app/internal/auth/domain"
)

type CreateAuthInput struct {
	ID          string
	ProfileID   string
	ProfileType string
	Email       string
	Password    string
	CreatedAt   time.Time
}

type CreateAuthUseCase struct {
	repo   domain.AuthRepository
	hasher domain.PasswordHasher
}

func NewCreateAuthUseCase(repo domain.AuthRepository, hasher domain.PasswordHasher) *CreateAuthUseCase {
	return &CreateAuthUseCase{repo: repo, hasher: hasher}
}

func (uc *CreateAuthUseCase) Execute(ctx context.Context, input CreateAuthInput) error {
	if input.Password == "" {
		return fmt.Errorf("la contraseña no puede estar vacía")
	}

	salt := input.CreatedAt.Format(time.RFC3339Nano)
	passwordHash := uc.hasher.Hash(input.Password, salt)

	authUser := &domain.AuthUser{
		ID:           input.ID,
		ProfileID:    input.ProfileID,
		ProfileType:  input.ProfileType,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.CreatedAt,
	}

	return uc.repo.Save(ctx, authUser)
}
