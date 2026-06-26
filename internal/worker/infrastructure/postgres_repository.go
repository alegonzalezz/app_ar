package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gcp-serverless-app/internal/worker/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, worker *domain.Worker) error {
	query := `
		INSERT INTO workers (id, name, email, phone, role, collective_agreement, salary, hire_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			role = EXCLUDED.role,
			collective_agreement = EXCLUDED.collective_agreement,
			salary = EXCLUDED.salary,
			hire_date = EXCLUDED.hire_date,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		worker.ID, worker.Name, worker.Email, worker.Phone,
		worker.Role, worker.CollectiveAgreement, worker.Salary, worker.HireDate,
		worker.CreatedAt, worker.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando trabajador: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Worker, error) {
	query := `
		SELECT id, name, email, phone, role, collective_agreement, salary, hire_date, created_at, updated_at
		FROM workers
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var w domain.Worker
	err := row.Scan(
		&w.ID, &w.Name, &w.Email, &w.Phone,
		&w.Role, &w.CollectiveAgreement, &w.Salary, &w.HireDate,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrWorkerNotFound
		}
		return nil, fmt.Errorf("error obteniendo trabajador por id: %w", err)
	}

	return &w, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*domain.Worker, error) {
	query := `
		SELECT id, name, email, phone, role, collective_agreement, salary, hire_date, created_at, updated_at
		FROM workers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando trabajadores: %w", err)
	}
	defer rows.Close()

	var workers []*domain.Worker
	for rows.Next() {
		var w domain.Worker
		if err := rows.Scan(
			&w.ID, &w.Name, &w.Email, &w.Phone,
			&w.Role, &w.CollectiveAgreement, &w.Salary, &w.HireDate,
			&w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando trabajador: %w", err)
		}
		workers = append(workers, &w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando trabajadores: %w", err)
	}

	return workers, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE workers
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error eliminando trabajador: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrWorkerNotFound
	}

	return nil
}

func (r *PostgresRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM workers WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando email de trabajador: %w", err)
	}

	return exists, nil
}
