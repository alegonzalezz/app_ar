package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	response "gcp-serverless-app/internal/shared/platform/http"
	taskApp "gcp-serverless-app/internal/task/application"
	taskDomain "gcp-serverless-app/internal/task/domain"
)

type TaskHandler struct {
	createUC *taskApp.CreateTaskUseCase
	getUC    *taskApp.GetTaskUseCase
	listUC   *taskApp.ListTasksUseCase
	updateUC *taskApp.UpdateTaskUseCase
	deleteUC *taskApp.DeleteTaskUseCase
	statusUC *taskApp.UpdateTaskStatusUseCase
}

func NewTaskHandler(
	createUC *taskApp.CreateTaskUseCase,
	getUC *taskApp.GetTaskUseCase,
	listUC *taskApp.ListTasksUseCase,
	updateUC *taskApp.UpdateTaskUseCase,
	deleteUC *taskApp.DeleteTaskUseCase,
	statusUC *taskApp.UpdateTaskStatusUseCase,
) *TaskHandler {
	return &TaskHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		statusUC: statusUC,
	}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /tasks", h.handleCreate)
	mux.HandleFunc("GET /tasks", h.handleList)
	mux.HandleFunc("GET /tasks/{id}", h.handleGet)
	mux.HandleFunc("PUT /tasks/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
	mux.HandleFunc("PATCH /tasks/{id}/status", h.handleUpdateStatus)
}

type taskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	Cost        *float64   `json:"cost"`
	CustomerID  string     `json:"customer_id"`
	WorkerID    string     `json:"worker_id"`
	DueDate     *time.Time `json:"due_date"`
}

type taskStatusRequest struct {
	Status string `json:"status"`
}

type taskDeleteRequest struct {
	Reason *string `json:"reason"`
}

func (h *TaskHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Priority = strings.TrimSpace(req.Priority)
	req.CustomerID = strings.TrimSpace(req.CustomerID)
	req.WorkerID = strings.TrimSpace(req.WorkerID)

	var errs []response.ErrorDetail
	if req.Title == "" {
		errs = append(errs, response.ErrorDetail{Field: "title", Code: "required_field"})
	}
	if req.CustomerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "customer_id", Code: "required_field"})
	}
	if req.WorkerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "worker_id", Code: "required_field"})
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	validPriority := false
	for _, p := range taskDomain.ValidPriorities {
		if p == req.Priority {
			validPriority = true
			break
		}
	}
	if !validPriority {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "priority", Code: "invalid_value"})
		return
	}

	task, err := h.createUC.Execute(r.Context(), taskApp.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Cost:        req.Cost,
		CustomerID:  req.CustomerID,
		WorkerID:    req.WorkerID,
		DueDate:     req.DueDate,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, task)
}

func (h *TaskHandler) handleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var customerID, workerID, status, priority *string
	var dueDateFrom, dueDateTo *time.Time

	if v := q.Get("customer_id"); v != "" {
		customerID = &v
	}
	if v := q.Get("worker_id"); v != "" {
		workerID = &v
	}
	if v := q.Get("status"); v != "" {
		status = &v
	}
	if v := q.Get("priority"); v != "" {
		priority = &v
	}
	if v := q.Get("due_date_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dueDateFrom = &t
		}
	}
	if v := q.Get("due_date_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dueDateTo = &t
		}
	}

	tasks, err := h.listUC.Execute(r.Context(), taskApp.ListTaskInput{
		CustomerID:  customerID,
		WorkerID:    workerID,
		Status:      status,
		Priority:    priority,
		DueDateFrom: dueDateFrom,
		DueDateTo:   dueDateTo,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if tasks == nil {
		tasks = []*taskDomain.Task{}
	}

	response.Success(w, http.StatusOK, tasks)
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	task, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == taskDomain.ErrTaskNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, task)
}

func (h *TaskHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Priority = strings.TrimSpace(req.Priority)
	req.CustomerID = strings.TrimSpace(req.CustomerID)
	req.WorkerID = strings.TrimSpace(req.WorkerID)

	var errs []response.ErrorDetail
	if req.Title == "" {
		errs = append(errs, response.ErrorDetail{Field: "title", Code: "required_field"})
	}
	if req.CustomerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "customer_id", Code: "required_field"})
	}
	if req.WorkerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "worker_id", Code: "required_field"})
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	validPriority := false
	for _, p := range taskDomain.ValidPriorities {
		if p == req.Priority {
			validPriority = true
			break
		}
	}
	if !validPriority {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "priority", Code: "invalid_value"})
		return
	}

	task, err := h.updateUC.Execute(r.Context(), taskApp.UpdateTaskInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Cost:        req.Cost,
		CustomerID:  req.CustomerID,
		WorkerID:    req.WorkerID,
		DueDate:     req.DueDate,
	})

	if err != nil {
		if err == taskDomain.ErrTaskNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, task)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req taskDeleteRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.deleteUC.Execute(r.Context(), taskApp.DeleteTaskInput{
		ID:     id,
		Reason: req.Reason,
	})

	if err != nil {
		if err == taskDomain.ErrTaskNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req taskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Status = strings.TrimSpace(req.Status)

	if req.Status == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "status", Code: "required_field"})
		return
	}

	task, err := h.statusUC.Execute(r.Context(), taskApp.UpdateTaskStatusInput{
		ID:     id,
		Status: req.Status,
	})

	if err != nil {
		if err == taskDomain.ErrTaskNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_found"})
			return
		}
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "status", Code: "invalid_value"})
		return
	}

	response.Success(w, http.StatusOK, task)
}
