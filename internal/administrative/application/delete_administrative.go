package application

import (
	"context"

	"gcp-serverless-app/internal/administrative/domain"
)

type DeleteAdministrativeUseCase struct {
	repo domain.AdministrativeRepository
}

func NewDeleteAdministrativeUseCase(repo domain.AdministrativeRepository) *DeleteAdministrativeUseCase {
	return &DeleteAdministrativeUseCase{repo: repo}
}

func (uc *DeleteAdministrativeUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
