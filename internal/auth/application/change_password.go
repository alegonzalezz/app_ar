package application

import (
	"context"

	domain "gcp-serverless-app/internal/auth/domain"
)

type ChangePasswordInput struct {
	Email       string
	OldPassword string
	NewPassword string
}

type ChangePasswordUseCase struct {
	repo   domain.AuthRepository
	hasher domain.PasswordHasher
}

func NewChangePasswordUseCase(repo domain.AuthRepository, hasher domain.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo, hasher: hasher}
}

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

	return uc.repo.UpdatePassword(ctx, authUser.ID, newHash)
}
