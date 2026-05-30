package application

import (
	"context"

	domain "gcp-serverless-app/internal/greeting/domain"
)

// FindGreetingUseCase busca un saludo por ID.
type FindGreetingUseCase struct {
	repo domain.GreetingRepository
}

// NewFindGreetingUseCase crea una nueva instancia del caso de uso.
func NewFindGreetingUseCase(repo domain.GreetingRepository) *FindGreetingUseCase {
	return &FindGreetingUseCase{repo: repo}
}

// Execute ejecuta la búsqueda del saludo.
func (uc *FindGreetingUseCase) Execute(ctx context.Context, id string) (*domain.Greeting, error) {
	if id == "" {
		id = "1"
	}
	return uc.repo.GetByID(ctx, id)
}
