package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	appointmentApp "gcp-serverless-app/internal/appointment/application"
	appointmentDomain "gcp-serverless-app/internal/appointment/domain"
	response "gcp-serverless-app/internal/shared/platform/http"
)

type AppointmentHandler struct {
	createUC                *appointmentApp.CreateAppointmentUseCase
	getUC                   *appointmentApp.GetAppointmentUseCase
	listUC                  *appointmentApp.ListAppointmentsUseCase
	updateUC                *appointmentApp.UpdateAppointmentUseCase
	deleteUC                *appointmentApp.DeleteAppointmentUseCase
	statusUC                *appointmentApp.UpdateAppointmentStatusUseCase
	assignTaskUC            *appointmentApp.AssignTaskUseCase
	unassignTaskUC          *appointmentApp.UnassignTaskUseCase
	getTasksByAppointmentUC *appointmentApp.GetTasksByAppointmentUseCase
	getAppointmentsByTaskUC *appointmentApp.GetAppointmentsByTaskUseCase
}

func NewAppointmentHandler(
	createUC *appointmentApp.CreateAppointmentUseCase,
	getUC *appointmentApp.GetAppointmentUseCase,
	listUC *appointmentApp.ListAppointmentsUseCase,
	updateUC *appointmentApp.UpdateAppointmentUseCase,
	deleteUC *appointmentApp.DeleteAppointmentUseCase,
	statusUC *appointmentApp.UpdateAppointmentStatusUseCase,
	assignTaskUC *appointmentApp.AssignTaskUseCase,
	unassignTaskUC *appointmentApp.UnassignTaskUseCase,
	getTasksByAppointmentUC *appointmentApp.GetTasksByAppointmentUseCase,
	getAppointmentsByTaskUC *appointmentApp.GetAppointmentsByTaskUseCase,
) *AppointmentHandler {
	return &AppointmentHandler{
		createUC:                createUC,
		getUC:                   getUC,
		listUC:                  listUC,
		updateUC:                updateUC,
		deleteUC:                deleteUC,
		statusUC:                statusUC,
		assignTaskUC:            assignTaskUC,
		unassignTaskUC:          unassignTaskUC,
		getTasksByAppointmentUC: getTasksByAppointmentUC,
		getAppointmentsByTaskUC: getAppointmentsByTaskUC,
	}
}

func (h *AppointmentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /appointments", h.handleCreate)
	mux.HandleFunc("GET /appointments", h.handleList)
	mux.HandleFunc("GET /appointments/{id}", h.handleGet)
	mux.HandleFunc("PUT /appointments/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /appointments/{id}", h.handleDelete)
	mux.HandleFunc("PATCH /appointments/{id}/status", h.handleUpdateStatus)
	mux.HandleFunc("POST /appointments/{id}/tasks", h.handleAssignTask)
	mux.HandleFunc("DELETE /appointments/{id}/tasks/{taskId}", h.handleUnassignTask)
	mux.HandleFunc("GET /appointments/{id}/tasks", h.handleGetTasksByAppointment)
	mux.HandleFunc("GET /tasks/{id}/appointments", h.handleGetAppointmentsByTask)
}

type appointmentRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	CustomerID  string     `json:"customer_id"`
	WorkerID    string     `json:"worker_id"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Notes       *string    `json:"notes"`
}

type appointmentStatusRequest struct {
	Status string `json:"status"`
}

type appointmentDeleteRequest struct {
	Reason *string `json:"reason"`
}

type assignTaskRequest struct {
	TaskID string  `json:"task_id"`
	Notes  *string `json:"notes"`
}

func (h *AppointmentHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req appointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
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
	if req.StartTime == nil {
		errs = append(errs, response.ErrorDetail{Field: "start_time", Code: "required_field"})
	}
	if req.EndTime == nil {
		errs = append(errs, response.ErrorDetail{Field: "end_time", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	appointment, err := h.createUC.Execute(r.Context(), appointmentApp.CreateAppointmentInput{
		Title:       req.Title,
		Description: req.Description,
		CustomerID:  req.CustomerID,
		WorkerID:    req.WorkerID,
		StartTime:   *req.StartTime,
		EndTime:     *req.EndTime,
		Notes:       req.Notes,
	})

	if err != nil {
		switch err {
		case appointmentDomain.ErrInvalidTimeRange:
			response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_time_range"})
		case appointmentDomain.ErrTimeConflict:
			response.Error(w, http.StatusConflict, response.ErrorDetail{Code: "time_conflict"})
		default:
			response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		}
		return
	}

	response.Success(w, http.StatusCreated, appointment)
}

func (h *AppointmentHandler) handleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var customerID, workerID, status *string
	var dateFrom, dateTo *time.Time

	if v := q.Get("customer_id"); v != "" {
		customerID = &v
	}
	if v := q.Get("worker_id"); v != "" {
		workerID = &v
	}
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

	appointments, err := h.listUC.Execute(r.Context(), appointmentApp.ListAppointmentsInput{
		CustomerID: customerID,
		WorkerID:   workerID,
		Status:     status,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if appointments == nil {
		appointments = []*appointmentDomain.Appointment{}
	}

	response.Success(w, http.StatusOK, appointments)
}

func (h *AppointmentHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	appointment, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == appointmentDomain.ErrAppointmentNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "appointment_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req appointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
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
	if req.StartTime == nil {
		errs = append(errs, response.ErrorDetail{Field: "start_time", Code: "required_field"})
	}
	if req.EndTime == nil {
		errs = append(errs, response.ErrorDetail{Field: "end_time", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	appointment, err := h.updateUC.Execute(r.Context(), appointmentApp.UpdateAppointmentInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		CustomerID:  req.CustomerID,
		WorkerID:    req.WorkerID,
		StartTime:   *req.StartTime,
		EndTime:     *req.EndTime,
		Notes:       req.Notes,
	})

	if err != nil {
		switch err {
		case appointmentDomain.ErrAppointmentNotFound:
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "appointment_not_found"})
		case appointmentDomain.ErrInvalidTimeRange:
			response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_time_range"})
		case appointmentDomain.ErrTimeConflict:
			response.Error(w, http.StatusConflict, response.ErrorDetail{Code: "time_conflict"})
		default:
			response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		}
		return
	}

	response.Success(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req appointmentDeleteRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.deleteUC.Execute(r.Context(), appointmentApp.DeleteAppointmentInput{
		ID:     id,
		Reason: req.Reason,
	})

	if err != nil {
		if err == appointmentDomain.ErrAppointmentNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "appointment_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AppointmentHandler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req appointmentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Status = strings.TrimSpace(req.Status)

	if req.Status == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "status", Code: "required_field"})
		return
	}

	appointment, err := h.statusUC.Execute(r.Context(), appointmentApp.UpdateAppointmentStatusInput{
		ID:     id,
		Status: req.Status,
	})

	if err != nil {
		if err == appointmentDomain.ErrAppointmentNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "appointment_not_found"})
			return
		}
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "status", Code: "invalid_value"})
		return
	}

	response.Success(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) handleAssignTask(w http.ResponseWriter, r *http.Request) {
	appointmentID := r.PathValue("id")
	if appointmentID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req assignTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.TaskID = strings.TrimSpace(req.TaskID)

	if req.TaskID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "task_id", Code: "required_field"})
		return
	}

	err := h.assignTaskUC.Execute(r.Context(), appointmentApp.AssignTaskInput{
		AppointmentID: appointmentID,
		TaskID:        req.TaskID,
		Notes:         req.Notes,
	})

	if err != nil {
		if err == appointmentDomain.ErrTaskAlreadyAssigned {
			response.Error(w, http.StatusConflict, response.ErrorDetail{Code: "already_assigned"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, map[string]string{"status": "assigned"})
}

func (h *AppointmentHandler) handleUnassignTask(w http.ResponseWriter, r *http.Request) {
	appointmentID := r.PathValue("id")
	taskID := r.PathValue("taskId")

	if appointmentID == "" || taskID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "required_field"})
		return
	}

	err := h.unassignTaskUC.Execute(r.Context(), appointmentApp.UnassignTaskInput{
		AppointmentID: appointmentID,
		TaskID:        taskID,
	})

	if err != nil {
		if err == appointmentDomain.ErrTaskNotAssigned {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_assigned"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"status": "unassigned"})
}

func (h *AppointmentHandler) handleGetTasksByAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentID := r.PathValue("id")
	if appointmentID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	tasks, err := h.getTasksByAppointmentUC.Execute(r.Context(), appointmentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if tasks == nil {
		tasks = []appointmentDomain.AppointmentTask{}
	}

	response.Success(w, http.StatusOK, tasks)
}

func (h *AppointmentHandler) handleGetAppointmentsByTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	appointments, err := h.getAppointmentsByTaskUC.Execute(r.Context(), taskID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if appointments == nil {
		appointments = []*appointmentDomain.Appointment{}
	}

	response.Success(w, http.StatusOK, appointments)
}
