package domain

import "time"

// Customer representa la entidad de cliente del dominio.
type Customer struct {
	ID               string
	Name             string
	PhoneNumber      string
	ExtraPhoneNumber *string
	ContactEmail     string
	ManagerName      string
	Address          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
