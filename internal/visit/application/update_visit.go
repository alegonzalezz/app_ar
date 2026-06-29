package application

import (
	"context"
	"fmt"
	"time"

	"gcp-serverless-app/internal/visit/domain"
)

type UpdateVisitInput struct {
	ID    string
	Notes *string
}

type UpdateVisitUseCase struct {
	repo domain.VisitRepository
}

func NewUpdateVisitUseCase(repo domain.VisitRepository) *UpdateVisitUseCase {
	return &UpdateVisitUseCase{repo: repo}
}

func (uc *UpdateVisitUseCase) Execute(ctx context.Context, input UpdateVisitInput) (*domain.Visit, error) {
	visit, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	visit.Notes = input.Notes
	visit.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, visit); err != nil {
		return nil, fmt.Errorf("error actualizando visita: %w", err)
	}

	return visit, nil
}
