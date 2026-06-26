package application

import (
	"context"

	"gcp-serverless-app/internal/worker/domain"
)

type GetWorkerUseCase struct {
	repo domain.WorkerRepository
}

func NewGetWorkerUseCase(repo domain.WorkerRepository) *GetWorkerUseCase {
	return &GetWorkerUseCase{repo: repo}
}

func (uc *GetWorkerUseCase) Execute(ctx context.Context, id string) (*domain.Worker, error) {
	return uc.repo.GetByID(ctx, id)
}
