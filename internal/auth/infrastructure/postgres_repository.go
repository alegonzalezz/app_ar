package infrastructure

import (
	"context"
	"database/sql"
	"strings"
	"time"

	domain "gcp-serverless-app/internal/auth/domain"
	"gcp-serverless-app/internal/shared/platform/postgres"
)

type postgresAuthRepository struct {
	db *sql.DB
}

// NewPostgresRepository crea una nueva instancia del repositorio PostgreSQL de auth.
func NewPostgresRepository(db *sql.DB) domain.AuthRepository {
	return &postgresAuthRepository{db: db}
}

func (r *postgresAuthRepository) getDB(ctx context.Context) postgres.SQLQueryer {
	return postgres.GetQueryer(ctx, r.db)
}

func (r *postgresAuthRepository) Save(ctx context.Context, authUser *domain.AuthUser) error {
	query := `
		INSERT INTO auth_users (user_id, email, password_hash, salt, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			salt = EXCLUDED.salt,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at
	`
	db := r.getDB(ctx)
	_, err := db.ExecContext(ctx, query,
		authUser.UserID,
		authUser.Email,
		authUser.PasswordHash,
		authUser.Salt,
		authUser.CreatedAt,
		authUser.UpdatedAt,
		authUser.DeletedAt,
	)
	return err
}

func (r *postgresAuthRepository) GetByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	query := `
		SELECT user_id, email, password_hash, salt, created_at, updated_at, deleted_at
		FROM auth_users
		WHERE email = $1 AND deleted_at IS NULL
	`
	db := r.getDB(ctx)
	var u domain.AuthUser
	err := db.QueryRowContext(ctx, query, strings.ToLower(email)).Scan(
		&u.UserID,
		&u.Email,
		&u.PasswordHash,
		&u.Salt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresAuthRepository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	query := `
		UPDATE auth_users
		SET password_hash = $1, updated_at = $2
		WHERE user_id = $3
	`
	db := r.getDB(ctx)
	_, err := db.ExecContext(ctx, query, newPasswordHash, time.Now().UTC(), userID)
	return err
}
