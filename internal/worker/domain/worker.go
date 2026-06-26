package domain

import "time"

type Worker struct {
	ID                  string
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	Salary              *float64
	HireDate            time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
