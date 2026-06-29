package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	domain "gcp-serverless-app/internal/user/domain"
)

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type CreateUserUseCase struct {
	txManager   domain.TxManager
	userRepo    domain.UserRepository
	authCreator domain.AuthCreator
}

func NewCreateUserUseCase(tx domain.TxManager, userRepo domain.UserRepository, authCreator domain.AuthCreator) *CreateUserUseCase {
	return &CreateUserUseCase{
		txManager:   tx,
		userRepo:    userRepo,
		authCreator: authCreator,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	var user *domain.User

	err := uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("error al verificar el email: %w", err)
		}
		if exists {
			return domain.ErrEmailAlreadyExists
		}

		now := time.Now().UTC()

		userID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("error al generar el id del usuario: %w", err)
		}

		authID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("error al generar el id de autenticación: %w", err)
		}

		user = &domain.User{
			ID:        userID,
			Email:     input.Email,
			Name:      input.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := uc.userRepo.Save(ctx, user); err != nil {
			return fmt.Errorf("error al guardar el usuario: %w", err)
		}

		if err := uc.authCreator.CreateAuth(ctx, authID, user.ID, "user", user.Email, input.Password, user.CreatedAt); err != nil {
			return fmt.Errorf("error al crear credenciales de autenticación: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
