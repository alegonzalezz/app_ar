package application

import (
	"context"

	"gcp-serverless-app/internal/appointment/domain"
)

type GetAppointmentUseCase struct {
	repo domain.AppointmentRepository
}

func NewGetAppointmentUseCase(repo domain.AppointmentRepository) *GetAppointmentUseCase {
	return &GetAppointmentUseCase{repo: repo}
}

func (uc *GetAppointmentUseCase) Execute(ctx context.Context, id string) (*domain.Appointment, error) {
	return uc.repo.GetByID(ctx, id)
}
