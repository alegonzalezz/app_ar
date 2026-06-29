package application

import (
	"context"
	"fmt"

	domain "gcp-serverless-app/internal/auth/domain"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	repo            domain.AuthRepository
	profileProvider domain.ProfileProvider
	hasher          domain.PasswordHasher
}

func NewLoginUseCase(repo domain.AuthRepository, profileProvider domain.ProfileProvider, hasher domain.PasswordHasher) *LoginUseCase {
	return &LoginUseCase{repo: repo, profileProvider: profileProvider, hasher: hasher}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*domain.UserInfo, error) {
	authUser, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	passwordHash := uc.hasher.Hash(input.Password, authUser.Salt)
	if passwordHash != authUser.PasswordHash {
		return nil, domain.ErrInvalidCredentials
	}

	profileName, profileData, err := uc.profileProvider.GetProfile(ctx, authUser.ProfileID, authUser.ProfileType)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo perfil: %w", err)
	}

	userInfo := &domain.UserInfo{
		ID:          authUser.ProfileID,
		Email:       authUser.Email,
		Name:        profileName,
		ProfileType: authUser.ProfileType,
		Profile:     profileData,
	}

	return userInfo, nil
}
