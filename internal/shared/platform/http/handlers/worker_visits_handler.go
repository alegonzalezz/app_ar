package handlers

import (
	"net/http"
	"strconv"
	"time"

	response "gcp-serverless-app/internal/shared/platform/http"
	visitApp "gcp-serverless-app/internal/visit/application"
	visitDomain "gcp-serverless-app/internal/visit/domain"
)

type WorkerVisitsHandler struct {
	listUC *visitApp.ListVisitsUseCase
}

func NewWorkerVisitsHandler(listUC *visitApp.ListVisitsUseCase) *WorkerVisitsHandler {
	return &WorkerVisitsHandler{listUC: listUC}
}

func (h *WorkerVisitsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /workers/{workerID}/visits", h.handleVisits)
}

func (h *WorkerVisitsHandler) handleVisits(w http.ResponseWriter, r *http.Request) {
	workerID := r.PathValue("workerID")
	if workerID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "workerID", Code: "required_field"})
		return
	}

	q := r.URL.Query()

	var dateFrom, dateTo *time.Time
	var status *string

	if v := q.Get("status"); v != "" {
		status = &v
	}
	if v := q.Get("date_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dateFrom = &t
		}
	}
	if v := q.Get("date_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dateTo = &t
		}
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.listUC.Execute(r.Context(), visitApp.ListVisitsInput{
		WorkerID: &workerID,
		Status:   status,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if result.Visits == nil {
		result.Visits = []*visitDomain.Visit{}
	}

	response.Success(w, http.StatusOK, result)
}
