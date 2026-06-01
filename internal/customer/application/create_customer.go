package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/customer/domain"
)

type CreateCustomerInput struct {
	Name             string
	PhoneNumber      string
	ExtraPhoneNumber *string
	ContactEmail     string
	ManagerName      string
	Address          string
}

type CreateCustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewCreateCustomerUseCase(repo domain.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{repo: repo}
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, input CreateCustomerInput) (*domain.Customer, error) {
	customerID, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de cliente: %w", err)
	}

	now := time.Now().UTC()
	customer := &domain.Customer{
		ID:               customerID,
		Name:             input.Name,
		PhoneNumber:      input.PhoneNumber,
		ExtraPhoneNumber: input.ExtraPhoneNumber,
		ContactEmail:     input.ContactEmail,
		ManagerName:      input.ManagerName,
		Address:          input.Address,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := uc.repo.Save(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
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
