package users

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gcp-serverless-app/auth"
)

// SQLQueryer define la interfaz común para sql.DB y sql.Tx
type SQLQueryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// postgresTxManager implementa TxManager usando transacciones SQL estándar
type postgresTxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &postgresTxManager{db: db}
}

func (m *postgresTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	// Propagar el *sql.Tx en el contexto usando la clave string "db_tx"
	txCtx := context.WithValue(ctx, "db_tx", tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error en la lógica (err: %v) y falló rollback (rbErr: %v)", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al hacer commit de la transacción: %w", err)
	}
	return nil
}

// postgresUserRepository implementa UserRepository
type postgresUserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) getDB(ctx context.Context) SQLQueryer {
	if tx, ok := ctx.Value("db_tx").(*sql.Tx); ok {
		return tx
	}
	return r.db
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

func (r *postgresUserRepository) Save(ctx context.Context, user *User) error {
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

// authBridge implementa la interfaz AuthCreator del dominio
type authBridge struct {
	authUseCase *auth.CreateAuthUseCase
}

func NewAuthBridge(auc *auth.CreateAuthUseCase) AuthCreator {
	return &authBridge{authUseCase: auc}
}

func (b *authBridge) CreateAuth(ctx context.Context, userID, email, password string, createdAt time.Time) error {
	return b.authUseCase.Execute(ctx, auth.CreateAuthInput{
		UserID:    userID,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	})
}

// CreateUserInput define los datos para el caso de uso
type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

// CreateUserUseCase coordina la creación de usuario y credenciales
type CreateUserUseCase struct {
	txManager   TxManager
	userRepo    UserRepository
	authCreator AuthCreator
}

func NewCreateUserUseCase(tx TxManager, userRepo UserRepository, authCreator AuthCreator) *CreateUserUseCase {
	return &CreateUserUseCase{
		txManager:   tx,
		userRepo:    userRepo,
		authCreator: authCreator,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*User, error) {
	var user *User

	err := uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		// 1. Validar que no exista el email
		exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("error al verificar el email: %w", err)
		}
		if exists {
			return ErrEmailAlreadyExists
		}

		// 2. Generar un nuevo ID de usuario (UUID v4)
		userID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("error al generar el id del usuario: %w", err)
		}

		now := time.Now().UTC()
		user = &User{
			ID:        userID,
			Email:     input.Email,
			Name:      input.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 3. Guardar el usuario
		if err := uc.userRepo.Save(ctx, user); err != nil {
			return fmt.Errorf("error al guardar el usuario: %w", err)
		}

		// 4. Crear la autenticación en la slice auth (transaccional)
		if err := uc.authCreator.CreateAuth(ctx, user.ID, user.Email, input.Password, user.CreatedAt); err != nil {
			return fmt.Errorf("error al crear credenciales de autenticación: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Handler maneja la petición HTTP para crear usuarios
type Handler struct {
	useCase *CreateUserUseCase
}

func NewHandler(uc *CreateUserUseCase) *Handler {
	return &Handler{useCase: uc}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido. Use POST."})
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "JSON malformado"})
		return
	}

	// Validaciones básicas
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Los campos 'name', 'email' y 'password' son requeridos"})
		return
	}

	user, err := h.useCase.Execute(r.Context(), CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if err == ErrEmailAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error interno del servidor", "details": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// generateUUID genera un UUID v4 usando la biblioteca estándar crypto/rand
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // versión 4
	b[8] = (b[8] & 0x3f) | 0x80 // variante RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
