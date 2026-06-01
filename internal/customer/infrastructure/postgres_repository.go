package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gcp-serverless-app/internal/customer/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (id, name, phone_number, extra_phone_number, contact_email, manager_name, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			phone_number = EXCLUDED.phone_number,
			extra_phone_number = EXCLUDED.extra_phone_number,
			contact_email = EXCLUDED.contact_email,
			manager_name = EXCLUDED.manager_name,
			address = EXCLUDED.address,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		customer.ID,
		customer.Name,
		customer.PhoneNumber,
		customer.ExtraPhoneNumber,
		customer.ContactEmail,
		customer.ManagerName,
		customer.Address,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error guardando cliente: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
		SELECT id, name, phone_number, extra_phone_number, contact_email, manager_name, address, created_at, updated_at
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var c domain.Customer
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.PhoneNumber,
		&c.ExtraPhoneNumber,
		&c.ContactEmail,
		&c.ManagerName,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("error obteniendo cliente por id: %w", err)
	}

	return &c, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*domain.Customer, error) {
	query := `
		SELECT id, name, phone_number, extra_phone_number, contact_email, manager_name, address, created_at, updated_at
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando clientes: %w", err)
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.PhoneNumber,
			&c.ExtraPhoneNumber,
			&c.ContactEmail,
			&c.ManagerName,
			&c.Address,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando cliente: %w", err)
		}
		customers = append(customers, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando clientes: %w", err)
	}

	return customers, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE customers
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error eliminando cliente: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrCustomerNotFound
	}

	return nil
}
