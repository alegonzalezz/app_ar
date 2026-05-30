package bridge

import (
	"context"

	authDomain "gcp-serverless-app/internal/auth/domain"
	userDomain "gcp-serverless-app/internal/user/domain"
)

// userProviderBridge implementa auth/domain.UserProvider
// conectando el módulo auth con el módulo user.
type userProviderBridge struct {
	userRepo userDomain.UserRepository
}

// NewUserProviderBridge crea un bridge que conecta auth con user para obtener datos del usuario.
func NewUserProviderBridge(repo userDomain.UserRepository) authDomain.UserProvider {
	return &userProviderBridge{userRepo: repo}
}

func (b *userProviderBridge) GetUser(ctx context.Context, userID string) (*authDomain.UserInfo, error) {
	user, err := b.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &authDomain.UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
