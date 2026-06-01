package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/customer/domain"
)

type UpdateCustomerInput struct {
	ID               string
	Name             string
	PhoneNumber      string
	ExtraPhoneNumber *string
	ContactEmail     string
	ManagerName      string
	Address          string
}

type UpdateCustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewUpdateCustomerUseCase(repo domain.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{repo: repo}
}

func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, input UpdateCustomerInput) (*domain.Customer, error) {
	customer, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	customer.Name = input.Name
	customer.PhoneNumber = input.PhoneNumber
	customer.ExtraPhoneNumber = input.ExtraPhoneNumber
	customer.ContactEmail = input.ContactEmail
	customer.ManagerName = input.ManagerName
	customer.Address = input.Address
	customer.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Save(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
