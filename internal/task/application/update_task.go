package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/task/domain"
)

type UpdateTaskInput struct {
	ID          string
	Title       string
	Description *string
	Priority    string
	Cost        *float64
	CustomerID  string
	WorkerID    string
	DueDate     *time.Time
}

type UpdateTaskUseCase struct {
	repo domain.TaskRepository
}

func NewUpdateTaskUseCase(repo domain.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input UpdateTaskInput) (*domain.Task, error) {
	task, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Priority = input.Priority
	task.Cost = input.Cost
	task.CustomerID = input.CustomerID
	task.WorkerID = input.WorkerID
	task.DueDate = input.DueDate
	task.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}
