package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAppointmentNotFound      = errors.New("turno no encontrado")
	ErrInvalidAppointmentStatus = errors.New("estado de turno inválido")
	ErrTimeConflict             = errors.New("el trabajador ya tiene un turno en ese horario")
	ErrInvalidTimeRange         = errors.New("end_time debe ser posterior a start_time")
	ErrTaskAlreadyAssigned      = errors.New("la tarea ya está asignada a este turno")
	ErrTaskNotAssigned          = errors.New("la tarea no está asignada a este turno")
)

var ValidAppointmentStatuses = []string{"scheduled", "confirmed", "in_progress", "completed", "cancelled"}

type AppointmentRepository interface {
	Save(ctx context.Context, appointment *Appointment) error
	GetByID(ctx context.Context, id string) (*Appointment, error)
	List(ctx context.Context, filter AppointmentFilter) ([]*Appointment, error)
	SoftDelete(ctx context.Context, id, reason string) error
	UpdateStatus(ctx context.Context, id, status string) error
	HasTimeConflict(ctx context.Context, workerID string, startTime, endTime time.Time, excludeID string) (bool, error)
	AssignTask(ctx context.Context, at *AppointmentTask) error
	UnassignTask(ctx context.Context, appointmentID, taskID string) error
	GetTasksByAppointment(ctx context.Context, appointmentID string) ([]AppointmentTask, error)
	GetAppointmentsByTask(ctx context.Context, taskID string) ([]*Appointment, error)
}

type AppointmentFilter struct {
	CustomerID *string
	WorkerID   *string
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
}
