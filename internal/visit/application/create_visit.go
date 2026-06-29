package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/visit/domain"
)

type CreateVisitInput struct {
	AppointmentID string
	CustomerID    string
	WorkerID      string
	Notes         *string
}

type CreateVisitUseCase struct {
	repo domain.VisitRepository
}

func NewCreateVisitUseCase(repo domain.VisitRepository) *CreateVisitUseCase {
	return &CreateVisitUseCase{repo: repo}
}

func (uc *CreateVisitUseCase) Execute(ctx context.Context, input CreateVisitInput) (*domain.Visit, error) {
	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de visita: %w", err)
	}

	now := time.Now().UTC()
	visit := &domain.Visit{
		ID:            id,
		AppointmentID: input.AppointmentID,
		CustomerID:    input.CustomerID,
		WorkerID:      input.WorkerID,
		Status:        "in_progress",
		Notes:         input.Notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.repo.Save(ctx, visit); err != nil {
		return nil, err
	}

	return visit, nil
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
