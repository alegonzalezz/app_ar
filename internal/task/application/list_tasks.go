package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/task/domain"
)

type ListTaskInput struct {
	CustomerID  *string
	WorkerID    *string
	Status      *string
	Priority    *string
	DueDateFrom *time.Time
	DueDateTo   *time.Time
}

type ListTasksUseCase struct {
	repo domain.TaskRepository
}

func NewListTasksUseCase(repo domain.TaskRepository) *ListTasksUseCase {
	return &ListTasksUseCase{repo: repo}
}

func (uc *ListTasksUseCase) Execute(ctx context.Context, input ListTaskInput) ([]*domain.Task, error) {
	filter := domain.TaskFilter{
		CustomerID:  input.CustomerID,
		WorkerID:    input.WorkerID,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDateFrom: input.DueDateFrom,
		DueDateTo:   input.DueDateTo,
	}

	return uc.repo.List(ctx, filter)
}
