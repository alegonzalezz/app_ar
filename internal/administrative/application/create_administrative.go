package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/administrative/domain"
)

type CreateAdministrativeInput struct {
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	WorkSchedule        string
	Salary              *float64
	HireDate            time.Time
}

type CreateAdministrativeUseCase struct {
	repo domain.AdministrativeRepository
}

func NewCreateAdministrativeUseCase(repo domain.AdministrativeRepository) *CreateAdministrativeUseCase {
	return &CreateAdministrativeUseCase{repo: repo}
}

func (uc *CreateAdministrativeUseCase) Execute(ctx context.Context, input CreateAdministrativeInput) (*domain.Administrative, error) {
	exists, err := uc.repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("error verificando email: %w", err)
	}
	if exists {
		return nil, domain.ErrAdministrativeEmailExists
	}

	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de administrativo: %w", err)
	}

	now := time.Now().UTC()
	admin := &domain.Administrative{
		ID:                  id,
		Name:                input.Name,
		Email:               input.Email,
		Phone:               input.Phone,
		Role:                input.Role,
		CollectiveAgreement: input.CollectiveAgreement,
		WorkSchedule:        input.WorkSchedule,
		HireDate:            input.HireDate,
		Salary:              input.Salary,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := uc.repo.Save(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
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
