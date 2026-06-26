package application

import (
	"context"

	"gcp-serverless-app/internal/task/domain"
)

type GetTaskUseCase struct {
	repo domain.TaskRepository
}

func NewGetTaskUseCase(repo domain.TaskRepository) *GetTaskUseCase {
	return &GetTaskUseCase{repo: repo}
}

func (uc *GetTaskUseCase) Execute(ctx context.Context, id string) (*domain.Task, error) {
	return uc.repo.GetByID(ctx, id)
}
