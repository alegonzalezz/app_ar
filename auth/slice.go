package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gcp-serverless-app/pkg/response"
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

func (r *postgresRepository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	query := `
		UPDATE auth_users
		SET password_hash = $1, updated_at = $2
		WHERE user_id = $3
	`
	db := r.getDB(ctx)
	_, err := db.ExecContext(ctx, query, newPasswordHash, time.Now().UTC(), userID)
	return err
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*AuthUser, error) {
	query := `
		SELECT user_id, email, password_hash, salt, created_at, updated_at, deleted_at
		FROM auth_users
		WHERE email = $1 AND deleted_at IS NULL
	`
	db := r.getDB(ctx)
	var u AuthUser
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
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
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

type ChangePasswordInput struct {
	Email       string
	OldPassword string
	NewPassword string
}

type ChangePasswordUseCase struct {
	repo Repository
}

func NewChangePasswordUseCase(repo Repository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
	authUser, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	currentHash := hashPassword(input.OldPassword, authUser.Salt)
	if currentHash != authUser.PasswordHash {
		return ErrInvalidCredentials
	}

	newHash := hashPassword(input.NewPassword, authUser.Salt)

	return uc.repo.UpdatePassword(ctx, authUser.UserID, newHash)
}

type ChangePasswordHandler struct {
	useCase *ChangePasswordUseCase
}

func NewChangePasswordHandler(uc *ChangePasswordUseCase) *ChangePasswordHandler {
	return &ChangePasswordHandler{useCase: uc}
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, response.ErrorDetail{Code: "method_not_allowed"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	var errs []response.ErrorDetail
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.OldPassword == "" {
		errs = append(errs, response.ErrorDetail{Field: "old_password", Code: "required_field"})
	}
	if req.NewPassword == "" {
		errs = append(errs, response.ErrorDetail{Field: "new_password", Code: "required_field"})
	} else if len(req.NewPassword) < 6 {
		errs = append(errs, response.ErrorDetail{Field: "new_password", Code: "min_length"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	err := h.useCase.Execute(r.Context(), ChangePasswordInput{
		Email:       req.Email,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		if err == ErrUserNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Field: "email", Code: "user_not_found"})
			return
		}
		if err == ErrInvalidCredentials {
			response.Error(w, http.StatusUnauthorized, response.ErrorDetail{Field: "old_password", Code: "invalid_user"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, "contraseña actualizada exitosamente")
}
