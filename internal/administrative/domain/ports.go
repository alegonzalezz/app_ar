package domain

import (
	"context"
	"errors"
)

var (
	ErrAdministrativeNotFound    = errors.New("administrativo no encontrado")
	ErrAdministrativeEmailExists = errors.New("el email del administrativo ya existe")
)

type AdministrativeRepository interface {
	Save(ctx context.Context, admin *Administrative) error
	GetByID(ctx context.Context, id string) (*Administrative, error)
	List(ctx context.Context) ([]*Administrative, error)
	SoftDelete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
