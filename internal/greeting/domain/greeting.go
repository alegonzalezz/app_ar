package domain

import "context"

// Greeting representa la entidad de saludo del dominio.
type Greeting struct {
	ID      string
	Message string
}

// GreetingRepository define el puerto driven para acceso a saludos.
type GreetingRepository interface {
	GetByID(ctx context.Context, id string) (*Greeting, error)
}
