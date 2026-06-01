package application

import (
	"context"

	"gcp-serverless-app/internal/customer/domain"
)

type ListCustomersUseCase struct {
	repo domain.CustomerRepository
}

func NewListCustomersUseCase(repo domain.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{repo: repo}
}

func (uc *ListCustomersUseCase) Execute(ctx context.Context) ([]*domain.Customer, error) {
	return uc.repo.List(ctx)
}
