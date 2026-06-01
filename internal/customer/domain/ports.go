package domain

import (
	"context"
	"errors"
)

var (
	ErrCustomerNotFound = errors.New("cliente no encontrado")
)

// CustomerRepository define el puerto driven para acceso a datos de clientes.
type CustomerRepository interface {
	Save(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id string) (*Customer, error)
	List(ctx context.Context) ([]*Customer, error)
	SoftDelete(ctx context.Context, id string) error
}
