package application

import (
	"context"

	"gcp-serverless-app/internal/visit/domain"
)

type UpdateVisitStatusInput struct {
	ID     string
	Status string
}

type UpdateVisitStatusUseCase struct {
	repo domain.VisitRepository
}

func NewUpdateVisitStatusUseCase(repo domain.VisitRepository) *UpdateVisitStatusUseCase {
	return &UpdateVisitStatusUseCase{repo: repo}
}

func (uc *UpdateVisitStatusUseCase) Execute(ctx context.Context, input UpdateVisitStatusInput) error {
	valid := false
	for _, s := range domain.ValidVisitStatuses {
		if s == input.Status {
			valid = true
			break
		}
	}
	if !valid {
		return domain.ErrInvalidVisitStatus
	}

	return uc.repo.UpdateStatus(ctx, input.ID, input.Status)
}

type AssignTaskInput struct {
	VisitID string
	TaskID  string
	Notes   *string
}

type AssignTaskUseCase struct {
	repo domain.VisitRepository
}

func NewAssignTaskUseCase(repo domain.VisitRepository) *AssignTaskUseCase {
	return &AssignTaskUseCase{repo: repo}
}

func (uc *AssignTaskUseCase) Execute(ctx context.Context, input AssignTaskInput) error {
	vt := &domain.VisitTask{
		VisitID: input.VisitID,
		TaskID:  input.TaskID,
		Notes:   input.Notes,
	}

	return uc.repo.AssignTask(ctx, vt)
}

type UnassignTaskInput struct {
	VisitID string
	TaskID  string
}

type UnassignTaskUseCase struct {
	repo domain.VisitRepository
}

func NewUnassignTaskUseCase(repo domain.VisitRepository) *UnassignTaskUseCase {
	return &UnassignTaskUseCase{repo: repo}
}

func (uc *UnassignTaskUseCase) Execute(ctx context.Context, input UnassignTaskInput) error {
	return uc.repo.UnassignTask(ctx, input.VisitID, input.TaskID)
}

type GetVisitTasksUseCase struct {
	repo domain.VisitRepository
}

func NewGetVisitTasksUseCase(repo domain.VisitRepository) *GetVisitTasksUseCase {
	return &GetVisitTasksUseCase{repo: repo}
}

func (uc *GetVisitTasksUseCase) Execute(ctx context.Context, visitID string) ([]domain.VisitTask, error) {
	_, err := uc.repo.GetByID(ctx, visitID)
	if err != nil {
		return nil, err
	}

	return uc.repo.GetTasksByVisit(ctx, visitID)
}
