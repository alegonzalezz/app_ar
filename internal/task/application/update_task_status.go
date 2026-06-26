package application

import (
	"context"
	"errors"

	"gcp-serverless-app/internal/task/domain"
)

type UpdateTaskStatusInput struct {
	ID     string
	Status string
}

type UpdateTaskStatusUseCase struct {
	repo domain.TaskRepository
}

func NewUpdateTaskStatusUseCase(repo domain.TaskRepository) *UpdateTaskStatusUseCase {
	return &UpdateTaskStatusUseCase{repo: repo}
}

func (uc *UpdateTaskStatusUseCase) Execute(ctx context.Context, input UpdateTaskStatusInput) (*domain.Task, error) {
	valid := false
	for _, s := range domain.ValidStatuses {
		if s == input.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("estado inválido: debe ser pending, in_progress, completed o cancelled")
	}

	if err := uc.repo.UpdateStatus(ctx, input.ID, input.Status); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, input.ID)
}
