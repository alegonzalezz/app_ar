package infrastructure

import (
	"context"
	"database/sql"
	"strings"

	"gcp-serverless-app/internal/shared/platform/postgres"
	domain "gcp-serverless-app/internal/user/domain"
)

type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgresRepository crea una nueva instancia del repositorio PostgreSQL de usuarios.
func NewPostgresRepository(db *sql.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) getDB(ctx context.Context) postgres.SQLQueryer {
	return postgres.GetQueryer(ctx, r.db)
}

func (r *postgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)"
	db := r.getDB(ctx)
	var exists bool
	err := db.QueryRowContext(ctx, query, strings.ToLower(email)).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := "SELECT id, email, name, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL"
	db := r.getDB(ctx)
	var u domain.User
	err := db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at
	`
	db := r.getDB(ctx)
	_, err := db.ExecContext(ctx, query,
		user.ID,
		strings.ToLower(user.Email),
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)
	return err
}
