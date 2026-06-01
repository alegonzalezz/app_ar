package application

import (
	"context"

	"gcp-serverless-app/internal/customer/domain"
)

type DeleteCustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewDeleteCustomerUseCase(repo domain.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{repo: repo}
}

func (uc *DeleteCustomerUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
