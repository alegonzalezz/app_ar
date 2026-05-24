package hello

import (
	"context"
	"database/sql"
	"net/http"

	"gcp-serverless-app/pkg/response"
)

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Greeting, error) {
	query := "SELECT id, message FROM greetings WHERE id = $1"
	var g Greeting
	err := r.db.QueryRowContext(ctx, query, id).Scan(&g.ID, &g.Message)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

type FindGreetingUseCase struct {
	repo Repository
}

func NewFindGreetingUseCase(repo Repository) *FindGreetingUseCase {
	return &FindGreetingUseCase{repo: repo}
}

func (uc *FindGreetingUseCase) Execute(ctx context.Context, id string) (*Greeting, error) {
	if id == "" {
		id = "1"
	}
	return uc.repo.GetByID(ctx, id)
}

type Handler struct {
	useCase *FindGreetingUseCase
}

func NewHandler(uc *FindGreetingUseCase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	result, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "not_found"})
		return
	}
	response.Success(w, http.StatusOK, result)
}
