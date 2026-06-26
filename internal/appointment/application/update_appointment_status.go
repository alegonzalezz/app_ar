package application

import (
	"context"
	"errors"

	"gcp-serverless-app/internal/appointment/domain"
)

type UpdateAppointmentStatusInput struct {
	ID     string
	Status string
}

type UpdateAppointmentStatusUseCase struct {
	repo domain.AppointmentRepository
}

func NewUpdateAppointmentStatusUseCase(repo domain.AppointmentRepository) *UpdateAppointmentStatusUseCase {
	return &UpdateAppointmentStatusUseCase{repo: repo}
}

func (uc *UpdateAppointmentStatusUseCase) Execute(ctx context.Context, input UpdateAppointmentStatusInput) (*domain.Appointment, error) {
	valid := false
	for _, s := range domain.ValidAppointmentStatuses {
		if s == input.Status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("estado inválido: debe ser scheduled, confirmed, in_progress, completed o cancelled")
	}

	if err := uc.repo.UpdateStatus(ctx, input.ID, input.Status); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, input.ID)
}
