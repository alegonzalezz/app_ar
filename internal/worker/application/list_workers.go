package application

import (
	"context"

	"gcp-serverless-app/internal/worker/domain"
)

type ListWorkersUseCase struct {
	repo domain.WorkerRepository
}

func NewListWorkersUseCase(repo domain.WorkerRepository) *ListWorkersUseCase {
	return &ListWorkersUseCase{repo: repo}
}

func (uc *ListWorkersUseCase) Execute(ctx context.Context) ([]*domain.Worker, error) {
	return uc.repo.List(ctx)
}
