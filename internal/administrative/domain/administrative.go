package domain

import "time"

type Administrative struct {
	ID                  string
	Name                string
	Email               string
	Phone               string
	Role                string
	CollectiveAgreement *string
	WorkSchedule        string
	HireDate            time.Time
	Salary              *float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
