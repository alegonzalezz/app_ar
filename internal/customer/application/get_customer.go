package application

import (
	"context"

	"gcp-serverless-app/internal/customer/domain"
)

type GetCustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewGetCustomerUseCase(repo domain.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{repo: repo}
}

func (uc *GetCustomerUseCase) Execute(ctx context.Context, id string) (*domain.Customer, error) {
	return uc.repo.GetByID(ctx, id)
}
