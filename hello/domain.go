package hello

import "context"

type Greeting struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Greeting, error)
}
