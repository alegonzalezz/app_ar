package domain

import "context"

type AuthRepository interface {
	Save(ctx context.Context, authUser *AuthUser) error
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
	UpdatePassword(ctx context.Context, authID, newPasswordHash string) error
}

type ProfileProvider interface {
	GetProfile(ctx context.Context, profileID string, profileType string) (profileName string, profileData interface{}, err error)
}

type PasswordHasher interface {
	Hash(password, salt string) string
}
