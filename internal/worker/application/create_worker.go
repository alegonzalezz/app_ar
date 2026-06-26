package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/worker/domain"
)

type CreateWorkerInput struct {
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	Salary              *float64
	HireDate            time.Time
}

type CreateWorkerUseCase struct {
	repo domain.WorkerRepository
}

func NewCreateWorkerUseCase(repo domain.WorkerRepository) *CreateWorkerUseCase {
	return &CreateWorkerUseCase{repo: repo}
}

func (uc *CreateWorkerUseCase) Execute(ctx context.Context, input CreateWorkerInput) (*domain.Worker, error) {
	exists, err := uc.repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("error verificando email: %w", err)
	}
	if exists {
		return nil, domain.ErrWorkerEmailExists
	}

	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de trabajador: %w", err)
	}

	now := time.Now().UTC()
	worker := &domain.Worker{
		ID:                  id,
		Name:                input.Name,
		Email:               input.Email,
		Phone:               input.Phone,
		Role:                input.Role,
		CollectiveAgreement: input.CollectiveAgreement,
		Salary:              input.Salary,
		HireDate:            input.HireDate,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := uc.repo.Save(ctx, worker); err != nil {
		return nil, err
	}

	return worker, nil
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
