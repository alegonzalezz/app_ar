package application

import (
	"context"

	"gcp-serverless-app/internal/appointment/domain"
)

type DeleteAppointmentInput struct {
	ID     string
	Reason *string
}

type DeleteAppointmentUseCase struct {
	repo domain.AppointmentRepository
}

func NewDeleteAppointmentUseCase(repo domain.AppointmentRepository) *DeleteAppointmentUseCase {
	return &DeleteAppointmentUseCase{repo: repo}
}

func (uc *DeleteAppointmentUseCase) Execute(ctx context.Context, input DeleteAppointmentInput) error {
	reason := ""
	if input.Reason != nil {
		reason = *input.Reason
	}
	return uc.repo.SoftDelete(ctx, input.ID, reason)
}
