package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gcp-serverless-app/internal/task/domain"
)

type CreateTaskInput struct {
	Title       string
	Description *string
	Priority    string
	Cost        *float64
	CustomerID  string
	WorkerID    string
	DueDate     *time.Time
}

type CreateTaskUseCase struct {
	repo domain.TaskRepository
}

func NewCreateTaskUseCase(repo domain.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{repo: repo}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskInput) (*domain.Task, error) {
	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generando id de tarea: %w", err)
	}

	now := time.Now().UTC()
	task := &domain.Task{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		Priority:    input.Priority,
		Cost:        input.Cost,
		CustomerID:  input.CustomerID,
		WorkerID:    input.WorkerID,
		DueDate:     input.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Save(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
