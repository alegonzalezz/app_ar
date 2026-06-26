package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/worker/domain"
)

type UpdateWorkerInput struct {
	ID                  string
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	Salary              *float64
	HireDate            time.Time
}

type UpdateWorkerUseCase struct {
	repo domain.WorkerRepository
}

func NewUpdateWorkerUseCase(repo domain.WorkerRepository) *UpdateWorkerUseCase {
	return &UpdateWorkerUseCase{repo: repo}
}

func (uc *UpdateWorkerUseCase) Execute(ctx context.Context, input UpdateWorkerInput) (*domain.Worker, error) {
	worker, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	worker.Name = input.Name
	worker.Email = input.Email
	worker.Phone = input.Phone
	worker.Role = input.Role
	worker.CollectiveAgreement = input.CollectiveAgreement
	worker.Salary = input.Salary
	worker.HireDate = input.HireDate
	worker.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, worker); err != nil {
		return nil, err
	}

	return worker, nil
}
