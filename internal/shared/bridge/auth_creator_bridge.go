package bridge

import (
	"context"
	"time"

	authApp "gcp-serverless-app/internal/auth/application"
	userDomain "gcp-serverless-app/internal/user/domain"
)

type authCreatorBridge struct {
	createAuthUC *authApp.CreateAuthUseCase
}

func NewAuthCreatorBridge(uc *authApp.CreateAuthUseCase) userDomain.AuthCreator {
	return &authCreatorBridge{createAuthUC: uc}
}

func (b *authCreatorBridge) CreateAuth(ctx context.Context, authID, profileID, profileType, email, password string, createdAt time.Time) error {
	return b.createAuthUC.Execute(ctx, authApp.CreateAuthInput{
		ID:          authID,
		ProfileID:   profileID,
		ProfileType: profileType,
		Email:       email,
		Password:    password,
		CreatedAt:   createdAt,
	})
}
