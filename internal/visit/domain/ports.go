package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrVisitNotFound       = errors.New("visita no encontrada")
	ErrInvalidVisitStatus  = errors.New("estado de visita inválido")
	ErrTaskAlreadyAssigned = errors.New("la tarea ya está asignada a esta visita")
	ErrTaskNotAssigned     = errors.New("la tarea no está asignada a esta visita")
	ErrAppointmentNotFound = errors.New("turno no encontrado")
)

var ValidVisitStatuses = []string{"in_progress", "completed", "cancelled"}

type VisitRepository interface {
	Save(ctx context.Context, visit *Visit) error
	GetByID(ctx context.Context, id string) (*Visit, error)
	List(ctx context.Context, filter VisitFilter) ([]*Visit, error)
	Count(ctx context.Context, filter VisitFilter) (int, error)
	SoftDelete(ctx context.Context, id, reason string) error
	UpdateStatus(ctx context.Context, id, status string) error
	AssignTask(ctx context.Context, vt *VisitTask) error
	UnassignTask(ctx context.Context, visitID, taskID string) error
	GetTasksByVisit(ctx context.Context, visitID string) ([]VisitTask, error)
}

type VisitFilter struct {
	CustomerID *string
	WorkerID   *string
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PageSize   int
}
