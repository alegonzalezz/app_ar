package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/appointment/domain"
)

type ListAppointmentsInput struct {
	CustomerID *string
	WorkerID   *string
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
}

type ListAppointmentsUseCase struct {
	repo domain.AppointmentRepository
}

func NewListAppointmentsUseCase(repo domain.AppointmentRepository) *ListAppointmentsUseCase {
	return &ListAppointmentsUseCase{repo: repo}
}

func (uc *ListAppointmentsUseCase) Execute(ctx context.Context, input ListAppointmentsInput) ([]*domain.Appointment, error) {
	filter := domain.AppointmentFilter{
		CustomerID: input.CustomerID,
		WorkerID:   input.WorkerID,
		Status:     input.Status,
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
	}

	return uc.repo.List(ctx, filter)
}
