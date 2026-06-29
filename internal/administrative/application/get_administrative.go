package application

import (
	"context"

	"gcp-serverless-app/internal/administrative/domain"
)

type GetAdministrativeUseCase struct {
	repo domain.AdministrativeRepository
}

func NewGetAdministrativeUseCase(repo domain.AdministrativeRepository) *GetAdministrativeUseCase {
	return &GetAdministrativeUseCase{repo: repo}
}

func (uc *GetAdministrativeUseCase) Execute(ctx context.Context, id string) (*domain.Administrative, error) {
	return uc.repo.GetByID(ctx, id)
}
