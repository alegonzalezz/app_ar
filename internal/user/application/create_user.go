package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	domain "gcp-serverless-app/internal/user/domain"
)

// CreateUserInput define los datos de entrada para crear un usuario.
type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

// CreateUserUseCase coordina la creación de usuario y credenciales.
type CreateUserUseCase struct {
	txManager   domain.TxManager
	userRepo    domain.UserRepository
	authCreator domain.AuthCreator
}

// NewCreateUserUseCase crea una nueva instancia del caso de uso.
func NewCreateUserUseCase(tx domain.TxManager, userRepo domain.UserRepository, authCreator domain.AuthCreator) *CreateUserUseCase {
	return &CreateUserUseCase{
		txManager:   tx,
		userRepo:    userRepo,
		authCreator: authCreator,
	}
}

// Execute ejecuta el caso de uso de creación de usuario.
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	var user *domain.User

	err := uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		// 1. Validar que no exista el email
		exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("error al verificar el email: %w", err)
		}
		if exists {
			return domain.ErrEmailAlreadyExists
		}

		// 2. Generar un nuevo ID de usuario (UUID v4)
		userID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("error al generar el id del usuario: %w", err)
		}

		now := time.Now().UTC()
		user = &domain.User{
			ID:        userID,
			Email:     input.Email,
			Name:      input.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 3. Guardar el usuario
		if err := uc.userRepo.Save(ctx, user); err != nil {
			return fmt.Errorf("error al guardar el usuario: %w", err)
		}

		// 4. Crear la autenticación (transaccional, vía bridge)
		if err := uc.authCreator.CreateAuth(ctx, user.ID, user.Email, input.Password, user.CreatedAt); err != nil {
			return fmt.Errorf("error al crear credenciales de autenticación: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// generateUUID genera un UUID v4 usando crypto/rand.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // versión 4
	b[8] = (b[8] & 0x3f) | 0x80 // variante RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
