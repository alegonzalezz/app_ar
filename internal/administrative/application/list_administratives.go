package application

import (
	"context"

	"gcp-serverless-app/internal/administrative/domain"
)

type ListAdministrativesUseCase struct {
	repo domain.AdministrativeRepository
}

func NewListAdministrativesUseCase(repo domain.AdministrativeRepository) *ListAdministrativesUseCase {
	return &ListAdministrativesUseCase{repo: repo}
}

func (uc *ListAdministrativesUseCase) Execute(ctx context.Context) ([]*domain.Administrative, error) {
	return uc.repo.List(ctx)
}
