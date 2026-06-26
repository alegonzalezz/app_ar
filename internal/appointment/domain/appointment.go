package domain

import "time"

type Appointment struct {
	ID              string
	Title           string
	Description     *string
	Status          string
	CustomerID      string
	WorkerID        string
	StartTime       time.Time
	EndTime         time.Time
	Notes           *string
	CancelledReason *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type AppointmentTask struct {
	AppointmentID string
	TaskID        string
	Notes         *string
}
