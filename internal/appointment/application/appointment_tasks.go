package application

import (
	"context"
	"fmt"

	"gcp-serverless-app/internal/appointment/domain"
)

type AssignTaskInput struct {
	AppointmentID string
	TaskID        string
	Notes         *string
}

type UnassignTaskInput struct {
	AppointmentID string
	TaskID        string
}

type AssignTaskUseCase struct {
	repo domain.AppointmentRepository
}

func NewAssignTaskUseCase(repo domain.AppointmentRepository) *AssignTaskUseCase {
	return &AssignTaskUseCase{repo: repo}
}

func (uc *AssignTaskUseCase) Execute(ctx context.Context, input AssignTaskInput) error {
	at := &domain.AppointmentTask{
		AppointmentID: input.AppointmentID,
		TaskID:        input.TaskID,
		Notes:         input.Notes,
	}

	if err := uc.repo.AssignTask(ctx, at); err != nil {
		return fmt.Errorf("error asignando tarea al turno: %w", err)
	}

	return nil
}

type UnassignTaskUseCase struct {
	repo domain.AppointmentRepository
}

func NewUnassignTaskUseCase(repo domain.AppointmentRepository) *UnassignTaskUseCase {
	return &UnassignTaskUseCase{repo: repo}
}

func (uc *UnassignTaskUseCase) Execute(ctx context.Context, input UnassignTaskInput) error {
	if err := uc.repo.UnassignTask(ctx, input.AppointmentID, input.TaskID); err != nil {
		return fmt.Errorf("error desasignando tarea del turno: %w", err)
	}

	return nil
}

type GetTasksByAppointmentUseCase struct {
	repo domain.AppointmentRepository
}

func NewGetTasksByAppointmentUseCase(repo domain.AppointmentRepository) *GetTasksByAppointmentUseCase {
	return &GetTasksByAppointmentUseCase{repo: repo}
}

func (uc *GetTasksByAppointmentUseCase) Execute(ctx context.Context, appointmentID string) ([]domain.AppointmentTask, error) {
	return uc.repo.GetTasksByAppointment(ctx, appointmentID)
}

type GetAppointmentsByTaskUseCase struct {
	repo domain.AppointmentRepository
}

func NewGetAppointmentsByTaskUseCase(repo domain.AppointmentRepository) *GetAppointmentsByTaskUseCase {
	return &GetAppointmentsByTaskUseCase{repo: repo}
}

func (uc *GetAppointmentsByTaskUseCase) Execute(ctx context.Context, taskID string) ([]*domain.Appointment, error) {
	return uc.repo.GetAppointmentsByTask(ctx, taskID)
}
