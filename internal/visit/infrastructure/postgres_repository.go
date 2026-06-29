package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gcp-serverless-app/internal/visit/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, v *domain.Visit) error {
	query := `
		INSERT INTO visits (id, appointment_id, customer_id, worker_id, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			customer_id = EXCLUDED.customer_id,
			worker_id = EXCLUDED.worker_id,
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		v.ID, v.AppointmentID, v.CustomerID, v.WorkerID,
		v.Status, v.Notes, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando visita: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Visit, error) {
	query := `
		SELECT id, appointment_id, customer_id, worker_id, status, notes, deleted_reason, created_at, updated_at
		FROM visits
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var v domain.Visit
	err := row.Scan(
		&v.ID, &v.AppointmentID, &v.CustomerID, &v.WorkerID,
		&v.Status, &v.Notes, &v.DeletedReason,
		&v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrVisitNotFound
		}
		return nil, fmt.Errorf("error obteniendo visita por id: %w", err)
	}

	return &v, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter domain.VisitFilter) ([]*domain.Visit, error) {
	query := `
		SELECT id, appointment_id, customer_id, worker_id, status, notes, deleted_reason, created_at, updated_at
		FROM visits
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
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIdx))
		args = append(args, *filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIdx))
		args = append(args, *filter.DateTo)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	page, pageSize := filter.Page, filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando visitas: %w", err)
	}
	defer rows.Close()

	var visits []*domain.Visit
	for rows.Next() {
		var v domain.Visit
		if err := rows.Scan(
			&v.ID, &v.AppointmentID, &v.CustomerID, &v.WorkerID,
			&v.Status, &v.Notes, &v.DeletedReason,
			&v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando visita: %w", err)
		}
		visits = append(visits, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando visitas: %w", err)
	}

	return visits, nil
}

func (r *PostgresRepository) Count(ctx context.Context, filter domain.VisitFilter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM visits
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

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando visitas: %w", err)
	}

	return count, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id, reason string) error {
	query := `
		UPDATE visits
		SET deleted_at = $1, deleted_reason = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), reason, id)
	if err != nil {
		return fmt.Errorf("error eliminando visita: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrVisitNotFound
	}

	return nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE visits
		SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error actualizando estado de visita: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrVisitNotFound
	}

	return nil
}

func (r *PostgresRepository) AssignTask(ctx context.Context, vt *domain.VisitTask) error {
	query := `
		INSERT INTO visit_tasks (visit_id, task_id, notes)
		VALUES ($1, $2, $3)
		ON CONFLICT (visit_id, task_id) DO NOTHING
	`

	res, err := r.db.ExecContext(ctx, query, vt.VisitID, vt.TaskID, vt.Notes)
	if err != nil {
		return fmt.Errorf("error asignando tarea a la visita: %w", err)
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

func (r *PostgresRepository) UnassignTask(ctx context.Context, visitID, taskID string) error {
	query := `
		DELETE FROM visit_tasks
		WHERE visit_id = $1 AND task_id = $2
	`

	res, err := r.db.ExecContext(ctx, query, visitID, taskID)
	if err != nil {
		return fmt.Errorf("error desasignando tarea de la visita: %w", err)
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

func (r *PostgresRepository) GetTasksByVisit(ctx context.Context, visitID string) ([]domain.VisitTask, error) {
	query := `
		SELECT visit_id, task_id, notes
		FROM visit_tasks
		WHERE visit_id = $1
		ORDER BY task_id
	`

	rows, err := r.db.QueryContext(ctx, query, visitID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tareas de la visita: %w", err)
	}
	defer rows.Close()

	var tasks []domain.VisitTask
	for rows.Next() {
		var vt domain.VisitTask
		if err := rows.Scan(&vt.VisitID, &vt.TaskID, &vt.Notes); err != nil {
			return nil, fmt.Errorf("error escaneando tarea de la visita: %w", err)
		}
		tasks = append(tasks, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando tareas de la visita: %w", err)
	}

	return tasks, nil
}
