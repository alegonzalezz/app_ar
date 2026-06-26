package application

import (
	"context"

	"gcp-serverless-app/internal/worker/domain"
)

type DeleteWorkerUseCase struct {
	repo domain.WorkerRepository
}

func NewDeleteWorkerUseCase(repo domain.WorkerRepository) *DeleteWorkerUseCase {
	return &DeleteWorkerUseCase{repo: repo}
}

func (uc *DeleteWorkerUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
