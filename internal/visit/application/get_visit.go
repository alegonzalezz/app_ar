package application

import (
	"context"

	"gcp-serverless-app/internal/visit/domain"
)

type GetVisitUseCase struct {
	repo domain.VisitRepository
}

func NewGetVisitUseCase(repo domain.VisitRepository) *GetVisitUseCase {
	return &GetVisitUseCase{repo: repo}
}

func (uc *GetVisitUseCase) Execute(ctx context.Context, id string) (*domain.Visit, error) {
	return uc.repo.GetByID(ctx, id)
}
