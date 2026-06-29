package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	response "gcp-serverless-app/internal/shared/platform/http"
	visitApp "gcp-serverless-app/internal/visit/application"
	visitDomain "gcp-serverless-app/internal/visit/domain"
)

type VisitHandler struct {
	createUC       *visitApp.CreateVisitUseCase
	getUC          *visitApp.GetVisitUseCase
	listUC         *visitApp.ListVisitsUseCase
	updateUC       *visitApp.UpdateVisitUseCase
	deleteUC       *visitApp.DeleteVisitUseCase
	statusUC       *visitApp.UpdateVisitStatusUseCase
	assignTaskUC   *visitApp.AssignTaskUseCase
	unassignTaskUC *visitApp.UnassignTaskUseCase
	getTasksUC     *visitApp.GetVisitTasksUseCase
}

func NewVisitHandler(
	createUC *visitApp.CreateVisitUseCase,
	getUC *visitApp.GetVisitUseCase,
	listUC *visitApp.ListVisitsUseCase,
	updateUC *visitApp.UpdateVisitUseCase,
	deleteUC *visitApp.DeleteVisitUseCase,
	statusUC *visitApp.UpdateVisitStatusUseCase,
	assignTaskUC *visitApp.AssignTaskUseCase,
	unassignTaskUC *visitApp.UnassignTaskUseCase,
	getTasksUC *visitApp.GetVisitTasksUseCase,
) *VisitHandler {
	return &VisitHandler{
		createUC:       createUC,
		getUC:          getUC,
		listUC:         listUC,
		updateUC:       updateUC,
		deleteUC:       deleteUC,
		statusUC:       statusUC,
		assignTaskUC:   assignTaskUC,
		unassignTaskUC: unassignTaskUC,
		getTasksUC:     getTasksUC,
	}
}

func (h *VisitHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /visits", h.handleCreate)
	mux.HandleFunc("GET /visits", h.handleList)
	mux.HandleFunc("GET /visits/{id}", h.handleGet)
	mux.HandleFunc("PUT /visits/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /visits/{id}", h.handleDelete)
	mux.HandleFunc("PATCH /visits/{id}/status", h.handleUpdateStatus)
	mux.HandleFunc("POST /visits/{id}/tasks", h.handleAssignTask)
	mux.HandleFunc("DELETE /visits/{id}/tasks/{taskId}", h.handleUnassignTask)
	mux.HandleFunc("GET /visits/{id}/tasks", h.handleGetTasks)
}

type visitRequest struct {
	AppointmentID string  `json:"appointment_id"`
	CustomerID    string  `json:"customer_id"`
	WorkerID      string  `json:"worker_id"`
	Notes         *string `json:"notes"`
}

type visitUpdateRequest struct {
	Notes *string `json:"notes"`
}

type visitStatusRequest struct {
	Status string `json:"status"`
}

type visitDeleteRequest struct {
	Reason string `json:"reason"`
}

type assignVisitTaskRequest struct {
	TaskID string  `json:"task_id"`
	Notes  *string `json:"notes"`
}

func (h *VisitHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req visitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.AppointmentID = strings.TrimSpace(req.AppointmentID)
	req.CustomerID = strings.TrimSpace(req.CustomerID)
	req.WorkerID = strings.TrimSpace(req.WorkerID)

	var errs []response.ErrorDetail
	if req.AppointmentID == "" {
		errs = append(errs, response.ErrorDetail{Field: "appointment_id", Code: "required_field"})
	}
	if req.CustomerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "customer_id", Code: "required_field"})
	}
	if req.WorkerID == "" {
		errs = append(errs, response.ErrorDetail{Field: "worker_id", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	visit, err := h.createUC.Execute(r.Context(), visitApp.CreateVisitInput{
		AppointmentID: req.AppointmentID,
		CustomerID:    req.CustomerID,
		WorkerID:      req.WorkerID,
		Notes:         req.Notes,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, visit)
}

func (h *VisitHandler) handleList(w http.ResponseWriter, r *http.Request) {
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

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.listUC.Execute(r.Context(), visitApp.ListVisitsInput{
		CustomerID: customerID,
		WorkerID:   workerID,
		Status:     status,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Page:       page,
		PageSize:   pageSize,
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

func (h *VisitHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	visit, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == visitDomain.ErrVisitNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "visit_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, visit)
}

func (h *VisitHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req visitUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	visit, err := h.updateUC.Execute(r.Context(), visitApp.UpdateVisitInput{
		ID:    id,
		Notes: req.Notes,
	})

	if err != nil {
		if err == visitDomain.ErrVisitNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "visit_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, visit)
}

func (h *VisitHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req visitDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Reason = strings.TrimSpace(req.Reason)

	if req.Reason == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "reason", Code: "required_field"})
		return
	}

	if err := h.deleteUC.Execute(r.Context(), id, req.Reason); err != nil {
		if err == visitDomain.ErrVisitNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "visit_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "visita eliminada correctamente"})
}

func (h *VisitHandler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req visitStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Status = strings.TrimSpace(req.Status)

	if req.Status == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "status", Code: "required_field"})
		return
	}

	if err := h.statusUC.Execute(r.Context(), visitApp.UpdateVisitStatusInput{
		ID:     id,
		Status: req.Status,
	}); err != nil {
		switch err {
		case visitDomain.ErrVisitNotFound:
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "visit_not_found"})
		case visitDomain.ErrInvalidVisitStatus:
			response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_status"})
		default:
			response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		}
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "estado actualizado correctamente"})
}

func (h *VisitHandler) handleAssignTask(w http.ResponseWriter, r *http.Request) {
	visitID := r.PathValue("id")
	if visitID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req assignVisitTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.TaskID = strings.TrimSpace(req.TaskID)

	if req.TaskID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "task_id", Code: "required_field"})
		return
	}

	if err := h.assignTaskUC.Execute(r.Context(), visitApp.AssignTaskInput{
		VisitID: visitID,
		TaskID:  req.TaskID,
		Notes:   req.Notes,
	}); err != nil {
		if err == visitDomain.ErrTaskAlreadyAssigned {
			response.Error(w, http.StatusConflict, response.ErrorDetail{Code: "task_already_assigned"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, map[string]string{"message": "tarea asignada correctamente"})
}

func (h *VisitHandler) handleUnassignTask(w http.ResponseWriter, r *http.Request) {
	visitID := r.PathValue("id")
	taskID := r.PathValue("taskId")

	if visitID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "taskId", Code: "required_field"})
		return
	}

	if err := h.unassignTaskUC.Execute(r.Context(), visitApp.UnassignTaskInput{
		VisitID: visitID,
		TaskID:  taskID,
	}); err != nil {
		if err == visitDomain.ErrTaskNotAssigned {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "task_not_assigned"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "tarea desasignada correctamente"})
}

func (h *VisitHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	visitID := r.PathValue("id")
	if visitID == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	tasks, err := h.getTasksUC.Execute(r.Context(), visitID)
	if err != nil {
		if err == visitDomain.ErrVisitNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "visit_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if tasks == nil {
		tasks = []visitDomain.VisitTask{}
	}

	response.Success(w, http.StatusOK, tasks)
}
