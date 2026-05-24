package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// SQLQueryer define la interfaz común para sql.DB y sql.Tx
type SQLQueryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *postgresRepository) getDB(ctx context.Context) SQLQueryer {
	// Buscamos la transacción en el contexto usando una clave string "db_tx"
	// para evitar acoplamiento e importaciones cíclicas entre slices.
	if tx, ok := ctx.Value("db_tx").(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *postgresRepository) Save(ctx context.Context, authUser *AuthUser) error {
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

type CreateAuthInput struct {
	UserID    string
	Email     string
	Password  string
	CreatedAt time.Time
}

type CreateAuthUseCase struct {
	repo Repository
}

func NewCreateAuthUseCase(repo Repository) *CreateAuthUseCase {
	return &CreateAuthUseCase{repo: repo}
}

func (uc *CreateAuthUseCase) Execute(ctx context.Context, input CreateAuthInput) error {
	if input.Password == "" {
		return fmt.Errorf("la contraseña no puede estar vacía")
	}

	// Generar el salt a partir de la fecha de creación en formato RFC3339Nano
	salt := input.CreatedAt.Format(time.RFC3339Nano)

	// Hashear la contraseña concatenada con el salt
	passwordHash := hashPassword(input.Password, salt)

	authUser := &AuthUser{
		UserID:       input.UserID,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.CreatedAt,
	}

	return uc.repo.Save(ctx, authUser)
}

func hashPassword(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
