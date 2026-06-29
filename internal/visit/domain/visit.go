package domain

import "time"

type Visit struct {
	ID            string
	AppointmentID string
	CustomerID    string
	WorkerID      string
	Status        string
	Notes         *string
	DeletedReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type VisitTask struct {
	VisitID string
	TaskID  string
	Notes   *string
}
