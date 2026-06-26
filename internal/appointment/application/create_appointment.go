package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/appointment/domain"
)

type CreateAppointmentInput struct {
	Title       string
	Description *string
	CustomerID  string
	WorkerID    string
	StartTime   time.Time
	EndTime     time.Time
	Notes       *string
}

type CreateAppointmentUseCase struct {
	repo domain.AppointmentRepository
}

func NewCreateAppointmentUseCase(repo domain.AppointmentRepository) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{repo: repo}
}

func (uc *CreateAppointmentUseCase) Execute(ctx context.Context, input CreateAppointmentInput) (*domain.Appointment, error) {
	if !input.EndTime.After(input.StartTime) {
		return nil, domain.ErrInvalidTimeRange
	}

	conflict, err := uc.repo.HasTimeConflict(ctx, input.WorkerID, input.StartTime, input.EndTime, "")
	if err != nil {
		return nil, fmt.Errorf("error verificando conflicto horario: %w", err)
	}
	if conflict {
		return nil, domain.ErrTimeConflict
	}

	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de turno: %w", err)
	}

	now := time.Now().UTC()
	appointment := &domain.Appointment{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Status:      "scheduled",
		CustomerID:  input.CustomerID,
		WorkerID:    input.WorkerID,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Notes:       input.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Save(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
