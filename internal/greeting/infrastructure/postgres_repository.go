package infrastructure

import (
	"context"
	"database/sql"

	domain "gcp-serverless-app/internal/greeting/domain"
)

type postgresGreetingRepository struct {
	db *sql.DB
}

// NewPostgresRepository crea una nueva instancia del repositorio PostgreSQL de saludos.
func NewPostgresRepository(db *sql.DB) domain.GreetingRepository {
	return &postgresGreetingRepository{db: db}
}

func (r *postgresGreetingRepository) GetByID(ctx context.Context, id string) (*domain.Greeting, error) {
	query := "SELECT id, message FROM greetings WHERE id = $1"
	var g domain.Greeting
	err := r.db.QueryRowContext(ctx, query, id).Scan(&g.ID, &g.Message)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
