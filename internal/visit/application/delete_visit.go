package application

import (
	"context"
	"fmt"

	"gcp-serverless-app/internal/visit/domain"
)

type DeleteVisitUseCase struct {
	repo domain.VisitRepository
}

func NewDeleteVisitUseCase(repo domain.VisitRepository) *DeleteVisitUseCase {
	return &DeleteVisitUseCase{repo: repo}
}

func (uc *DeleteVisitUseCase) Execute(ctx context.Context, id, reason string) error {
	if reason == "" {
		return fmt.Errorf("motivo de eliminación requerido")
	}

	return uc.repo.SoftDelete(ctx, id, reason)
}
