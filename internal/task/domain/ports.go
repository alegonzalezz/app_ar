package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTaskNotFound    = errors.New("tarea no encontrada")
	ErrInvalidStatus   = errors.New("estado inválido")
	ErrInvalidPriority = errors.New("prioridad inválida")
)

var ValidStatuses = []string{"pending", "in_progress", "completed", "cancelled"}
var ValidPriorities = []string{"low", "medium", "high", "critical"}

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]*Task, error)
	SoftDelete(ctx context.Context, id, reason string) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type TaskFilter struct {
	CustomerID  *string
	WorkerID    *string
	Status      *string
	Priority    *string
	DueDateFrom *time.Time
	DueDateTo   *time.Time
}
