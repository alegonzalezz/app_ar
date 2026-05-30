package bridge

import (
	"context"
	"time"

	authApp "gcp-serverless-app/internal/auth/application"
	userDomain "gcp-serverless-app/internal/user/domain"
)

// authCreatorBridge implementa user/domain.AuthCreator
// conectando el módulo user con el módulo auth.
type authCreatorBridge struct {
	createAuthUC *authApp.CreateAuthUseCase
}

// NewAuthCreatorBridge crea un bridge que conecta user con auth para crear credenciales.
func NewAuthCreatorBridge(uc *authApp.CreateAuthUseCase) userDomain.AuthCreator {
	return &authCreatorBridge{createAuthUC: uc}
}

func (b *authCreatorBridge) CreateAuth(ctx context.Context, userID, email, password string, createdAt time.Time) error {
	return b.createAuthUC.Execute(ctx, authApp.CreateAuthInput{
		UserID:    userID,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	})
}
