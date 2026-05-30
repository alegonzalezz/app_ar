package infrastructure

import (
	"crypto/sha256"
	"fmt"

	domain "gcp-serverless-app/internal/auth/domain"
)

type sha256Hasher struct{}

// NewSHA256Hasher crea una nueva instancia del hasher SHA256.
func NewSHA256Hasher() domain.PasswordHasher {
	return &sha256Hasher{}
}

// Hash hashea una contraseña con un salt usando SHA256.
func (h *sha256Hasher) Hash(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
