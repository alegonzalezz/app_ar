package domain

import "time"

type Task struct {
	ID            string
	Title         string
	Description   *string
	Status        string
	Priority      string
	Cost          *float64
	CustomerID    string
	WorkerID      string
	DueDate       *time.Time
	DeletedReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
