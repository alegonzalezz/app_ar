package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/appointment/domain"
)

type UpdateAppointmentInput struct {
	ID          string
	Title       string
	Description *string
	CustomerID  string
	WorkerID    string
	StartTime   time.Time
	EndTime     time.Time
	Notes       *string
}

type UpdateAppointmentUseCase struct {
	repo domain.AppointmentRepository
}

func NewUpdateAppointmentUseCase(repo domain.AppointmentRepository) *UpdateAppointmentUseCase {
	return &UpdateAppointmentUseCase{repo: repo}
}

func (uc *UpdateAppointmentUseCase) Execute(ctx context.Context, input UpdateAppointmentInput) (*domain.Appointment, error) {
	if !input.EndTime.After(input.StartTime) {
		return nil, domain.ErrInvalidTimeRange
	}

	appointment, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	conflict, err := uc.repo.HasTimeConflict(ctx, input.WorkerID, input.StartTime, input.EndTime, input.ID)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, domain.ErrTimeConflict
	}

	appointment.Title = input.Title
	appointment.Description = input.Description
	appointment.CustomerID = input.CustomerID
	appointment.WorkerID = input.WorkerID
	appointment.StartTime = input.StartTime
	appointment.EndTime = input.EndTime
	appointment.Notes = input.Notes
	appointment.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}
