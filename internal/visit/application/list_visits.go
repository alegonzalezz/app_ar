package application

import (
	"context"
	"time"

	"gcp-serverless-app/internal/visit/domain"
)

type ListVisitsInput struct {
	CustomerID *string
	WorkerID   *string
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PageSize   int
}

type ListVisitsResult struct {
	Visits   []*domain.Visit
	Total    int
	Page     int
	PageSize int
}

type ListVisitsUseCase struct {
	repo domain.VisitRepository
}

func NewListVisitsUseCase(repo domain.VisitRepository) *ListVisitsUseCase {
	return &ListVisitsUseCase{repo: repo}
}

func (uc *ListVisitsUseCase) Execute(ctx context.Context, input ListVisitsInput) (*ListVisitsResult, error) {
	filter := domain.VisitFilter{
		CustomerID: input.CustomerID,
		WorkerID:   input.WorkerID,
		Status:     input.Status,
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
		Page:       input.Page,
		PageSize:   input.PageSize,
	}

	total, err := uc.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	visits, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	page, pageSize := input.Page, input.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return &ListVisitsResult{
		Visits:   visits,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
