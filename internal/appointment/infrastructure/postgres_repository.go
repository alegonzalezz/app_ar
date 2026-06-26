package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gcp-serverless-app/internal/appointment/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, a *domain.Appointment) error {
	query := `
		INSERT INTO appointments (id, title, description, status, customer_id, worker_id, start_time, end_time, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			customer_id = EXCLUDED.customer_id,
			worker_id = EXCLUDED.worker_id,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		a.ID, a.Title, a.Description, a.Status, a.CustomerID,
		a.WorkerID, a.StartTime, a.EndTime, a.Notes,
		a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando turno: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	query := `
		SELECT id, title, description, status, customer_id, worker_id, start_time, end_time, notes, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var a domain.Appointment
	err := row.Scan(
		&a.ID, &a.Title, &a.Description, &a.Status, &a.CustomerID,
		&a.WorkerID, &a.StartTime, &a.EndTime, &a.Notes,
		&a.CancelledReason, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("error obteniendo turno por id: %w", err)
	}

	return &a, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter domain.AppointmentFilter) ([]*domain.Appointment, error) {
	query := `
		SELECT id, title, description, status, customer_id, worker_id, start_time, end_time, notes, cancelled_reason, created_at, updated_at
		FROM appointments
		WHERE deleted_at IS NULL
	`

	var args []interface{}
	argIdx := 1
	var conditions []string

	if filter.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argIdx))
		args = append(args, *filter.CustomerID)
		argIdx++
	}
	if filter.WorkerID != nil {
		conditions = append(conditions, fmt.Sprintf("worker_id = $%d", argIdx))
		args = append(args, *filter.WorkerID)
		argIdx++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", argIdx))
		args = append(args, *filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("start_time <= $%d", argIdx))
		args = append(args, *filter.DateTo)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY start_time ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando turnos: %w", err)
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(
			&a.ID, &a.Title, &a.Description, &a.Status, &a.CustomerID,
			&a.WorkerID, &a.StartTime, &a.EndTime, &a.Notes,
			&a.CancelledReason, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando turno: %w", err)
		}
		appointments = append(appointments, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando turnos: %w", err)
	}

	return appointments, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id, reason string) error {
	query := `
		UPDATE appointments
		SET deleted_at = $1, cancelled_reason = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), reason, id)
	if err != nil {
		return fmt.Errorf("error cancelando turno: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrAppointmentNotFound
	}

	return nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE appointments
		SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error actualizando estado de turno: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrAppointmentNotFound
	}

	return nil
}

func (r *PostgresRepository) HasTimeConflict(ctx context.Context, workerID string, startTime, endTime time.Time, excludeID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM appointments
		WHERE worker_id = $1
		AND deleted_at IS NULL
		AND start_time < $3
		AND end_time > $2
	`

	var args []interface{}
	args = append(args, workerID, startTime, endTime)
	argIdx := 4

	if excludeID != "" {
		query += fmt.Sprintf(" AND id != $%d", argIdx)
		args = append(args, excludeID)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error verificando conflicto horario: %w", err)
	}

	return count > 0, nil
}

func (r *PostgresRepository) AssignTask(ctx context.Context, at *domain.AppointmentTask) error {
	query := `
		INSERT INTO appointment_tasks (appointment_id, task_id, notes)
		VALUES ($1, $2, $3)
		ON CONFLICT (appointment_id, task_id) DO NOTHING
	`

	res, err := r.db.ExecContext(ctx, query, at.AppointmentID, at.TaskID, at.Notes)
	if err != nil {
		return fmt.Errorf("error asignando tarea al turno: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskAlreadyAssigned
	}

	return nil
}

func (r *PostgresRepository) UnassignTask(ctx context.Context, appointmentID, taskID string) error {
	query := `
		DELETE FROM appointment_tasks
		WHERE appointment_id = $1 AND task_id = $2
	`

	res, err := r.db.ExecContext(ctx, query, appointmentID, taskID)
	if err != nil {
		return fmt.Errorf("error desasignando tarea del turno: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotAssigned
	}

	return nil
}

func (r *PostgresRepository) GetTasksByAppointment(ctx context.Context, appointmentID string) ([]domain.AppointmentTask, error) {
	query := `
		SELECT appointment_id, task_id, notes
		FROM appointment_tasks
		WHERE appointment_id = $1
		ORDER BY task_id
	`

	rows, err := r.db.QueryContext(ctx, query, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tareas del turno: %w", err)
	}
	defer rows.Close()

	var tasks []domain.AppointmentTask
	for rows.Next() {
		var at domain.AppointmentTask
		if err := rows.Scan(&at.AppointmentID, &at.TaskID, &at.Notes); err != nil {
			return nil, fmt.Errorf("error escaneando tarea del turno: %w", err)
		}
		tasks = append(tasks, at)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando tareas del turno: %w", err)
	}

	return tasks, nil
}

func (r *PostgresRepository) GetAppointmentsByTask(ctx context.Context, taskID string) ([]*domain.Appointment, error) {
	query := `
		SELECT a.id, a.title, a.description, a.status, a.customer_id, a.worker_id, a.start_time, a.end_time, a.notes, a.cancelled_reason, a.created_at, a.updated_at
		FROM appointments a
		INNER JOIN appointment_tasks at ON at.appointment_id = a.id
		WHERE at.task_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.start_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo turnos de la tarea: %w", err)
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(
			&a.ID, &a.Title, &a.Description, &a.Status, &a.CustomerID,
			&a.WorkerID, &a.StartTime, &a.EndTime, &a.Notes,
			&a.CancelledReason, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando turno: %w", err)
		}
		appointments = append(appointments, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando turnos de la tarea: %w", err)
	}

	return appointments, nil
}
