package domain

import (
	"context"
	"errors"
)

var (
	ErrWorkerNotFound    = errors.New("trabajador no encontrado")
	ErrWorkerEmailExists = errors.New("el email del trabajador ya existe")
)

type WorkerRepository interface {
	Save(ctx context.Context, worker *Worker) error
	GetByID(ctx context.Context, id string) (*Worker, error)
	List(ctx context.Context) ([]*Worker, error)
	SoftDelete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
