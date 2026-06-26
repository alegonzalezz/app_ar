package application

import (
	"context"

	"gcp-serverless-app/internal/task/domain"
)

type DeleteTaskInput struct {
	ID     string
	Reason *string
}

type DeleteTaskUseCase struct {
	repo domain.TaskRepository
}

func NewDeleteTaskUseCase(repo domain.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{repo: repo}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, input DeleteTaskInput) error {
	reason := ""
	if input.Reason != nil {
		reason = *input.Reason
	}
	return uc.repo.SoftDelete(ctx, input.ID, reason)
}
