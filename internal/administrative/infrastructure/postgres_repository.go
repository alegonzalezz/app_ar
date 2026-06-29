package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gcp-serverless-app/internal/administrative/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, admin *domain.Administrative) error {
	query := `
		INSERT INTO administratives (id, name, email, phone, role, collective_agreement, work_schedule, hire_date, salary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			role = EXCLUDED.role,
			collective_agreement = EXCLUDED.collective_agreement,
			work_schedule = EXCLUDED.work_schedule,
			hire_date = EXCLUDED.hire_date,
			salary = EXCLUDED.salary,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx, query,
		admin.ID, admin.Name, admin.Email, admin.Phone,
		admin.Role, admin.CollectiveAgreement, admin.WorkSchedule, admin.HireDate,
		admin.Salary, admin.CreatedAt, admin.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando administrativo: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Administrative, error) {
	query := `
		SELECT id, name, email, phone, role, collective_agreement, work_schedule, hire_date, salary, created_at, updated_at
		FROM administratives
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var a domain.Administrative
	err := row.Scan(
		&a.ID, &a.Name, &a.Email, &a.Phone,
		&a.Role, &a.CollectiveAgreement, &a.WorkSchedule, &a.HireDate,
		&a.Salary, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAdministrativeNotFound
		}
		return nil, fmt.Errorf("error obteniendo administrativo por id: %w", err)
	}

	return &a, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*domain.Administrative, error) {
	query := `
		SELECT id, name, email, phone, role, collective_agreement, work_schedule, hire_date, salary, created_at, updated_at
		FROM administratives
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando administrativos: %w", err)
	}
	defer rows.Close()

	var administratives []*domain.Administrative
	for rows.Next() {
		var a domain.Administrative
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Email, &a.Phone,
			&a.Role, &a.CollectiveAgreement, &a.WorkSchedule, &a.HireDate,
			&a.Salary, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando administrativo: %w", err)
		}
		administratives = append(administratives, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando administrativos: %w", err)
	}

	return administratives, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE administratives
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error eliminando administrativo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrAdministrativeNotFound
	}

	return nil
}

func (r *PostgresRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM administratives WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando email de administrativo: %w", err)
	}

	return exists, nil
}
