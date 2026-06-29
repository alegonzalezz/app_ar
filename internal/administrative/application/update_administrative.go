package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/administrative/domain"
)

type UpdateAdministrativeInput struct {
	ID                  string
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	WorkSchedule        string
	Salary              *float64
	HireDate            time.Time
}

type UpdateAdministrativeUseCase struct {
	repo domain.AdministrativeRepository
}

func NewUpdateAdministrativeUseCase(repo domain.AdministrativeRepository) *UpdateAdministrativeUseCase {
	return &UpdateAdministrativeUseCase{repo: repo}
}

func (uc *UpdateAdministrativeUseCase) Execute(ctx context.Context, input UpdateAdministrativeInput) (*domain.Administrative, error) {
	admin, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	admin.Name = input.Name
	admin.Email = input.Email
	admin.Phone = input.Phone
	admin.Role = input.Role
	admin.CollectiveAgreement = input.CollectiveAgreement
	admin.WorkSchedule = input.WorkSchedule
	admin.Salary = input.Salary
	admin.HireDate = input.HireDate
	admin.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}
