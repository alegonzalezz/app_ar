package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gcp-serverless-app/internal/task/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, cost, customer_id, worker_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			priority = EXCLUDED.priority,
			cost = EXCLUDED.cost,
			customer_id = EXCLUDED.customer_id,
			worker_id = EXCLUDED.worker_id,
			due_date = EXCLUDED.due_date,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		task.ID, task.Title, task.Description, task.Status,
		task.Priority, task.Cost, task.CustomerID, task.WorkerID,
		task.DueDate, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando tarea: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, cost, customer_id, worker_id, due_date, deleted_reason, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var t domain.Task
	err := row.Scan(
		&t.ID, &t.Title, &t.Description, &t.Status,
		&t.Priority, &t.Cost, &t.CustomerID, &t.WorkerID,
		&t.DueDate, &t.DeletedReason, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("error obteniendo tarea por id: %w", err)
	}

	return &t, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, cost, customer_id, worker_id, due_date, deleted_reason, created_at, updated_at
		FROM tasks
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
	if filter.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, *filter.Priority)
		argIdx++
	}
	if filter.DueDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argIdx))
		args = append(args, *filter.DueDateFrom)
		argIdx++
	}
	if filter.DueDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argIdx))
		args = append(args, *filter.DueDateTo)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando tareas: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status,
			&t.Priority, &t.Cost, &t.CustomerID, &t.WorkerID,
			&t.DueDate, &t.DeletedReason, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando tarea: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando tareas: %w", err)
	}

	return tasks, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id, reason string) error {
	query := `
		UPDATE tasks
		SET deleted_at = $1, deleted_reason = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), reason, id)
	if err != nil {
		return fmt.Errorf("error eliminando tarea: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE tasks
		SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error actualizando estado de tarea: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
