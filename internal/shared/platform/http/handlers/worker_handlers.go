package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	response "gcp-serverless-app/internal/shared/platform/http"
	workerApp "gcp-serverless-app/internal/worker/application"
	workerDomain "gcp-serverless-app/internal/worker/domain"
)

type WorkerHandler struct {
	createUC *workerApp.CreateWorkerUseCase
	getUC    *workerApp.GetWorkerUseCase
	listUC   *workerApp.ListWorkersUseCase
	updateUC *workerApp.UpdateWorkerUseCase
	deleteUC *workerApp.DeleteWorkerUseCase
}

func NewWorkerHandler(
	createUC *workerApp.CreateWorkerUseCase,
	getUC *workerApp.GetWorkerUseCase,
	listUC *workerApp.ListWorkersUseCase,
	updateUC *workerApp.UpdateWorkerUseCase,
	deleteUC *workerApp.DeleteWorkerUseCase,
) *WorkerHandler {
	return &WorkerHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *WorkerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /workers", h.handleCreate)
	mux.HandleFunc("GET /workers", h.handleList)
	mux.HandleFunc("GET /workers/{id}", h.handleGet)
	mux.HandleFunc("PUT /workers/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /workers/{id}", h.handleDelete)
}

type workerRequest struct {
	Name                string   `json:"name"`
	Email               string   `json:"email"`
	Phone               string   `json:"phone"`
	Role                string   `json:"role"`
	CollectiveAgreement *string  `json:"collective_agreement"`
	Salary              *float64 `json:"salary"`
	HireDate            string   `json:"hire_date"`
}

func (h *WorkerHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req workerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Role = strings.TrimSpace(req.Role)
	req.HireDate = strings.TrimSpace(req.HireDate)

	var errs []response.ErrorDetail
	if req.Name == "" {
		errs = append(errs, response.ErrorDetail{Field: "name", Code: "required_field"})
	}
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.Phone == "" {
		errs = append(errs, response.ErrorDetail{Field: "phone", Code: "required_field"})
	}
	if req.Role == "" {
		errs = append(errs, response.ErrorDetail{Field: "role", Code: "required_field"})
	}
	if req.HireDate == "" {
		errs = append(errs, response.ErrorDetail{Field: "hire_date", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "hire_date", Code: "invalid_date"})
		return
	}

	worker, err := h.createUC.Execute(r.Context(), workerApp.CreateWorkerInput{
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Role:                req.Role,
		CollectiveAgreement: req.CollectiveAgreement,
		Salary:              req.Salary,
		HireDate:            hireDate,
	})

	if err != nil {
		if err == workerDomain.ErrWorkerEmailExists {
			response.Error(w, http.StatusConflict, response.ErrorDetail{Field: "email", Code: "already_exists"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, worker)
}

func (h *WorkerHandler) handleList(w http.ResponseWriter, r *http.Request) {
	workers, err := h.listUC.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if workers == nil {
		workers = []*workerDomain.Worker{}
	}

	response.Success(w, http.StatusOK, workers)
}

func (h *WorkerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	worker, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == workerDomain.ErrWorkerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "worker_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, worker)
}

func (h *WorkerHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req workerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Role = strings.TrimSpace(req.Role)
	req.HireDate = strings.TrimSpace(req.HireDate)

	var errs []response.ErrorDetail
	if req.Name == "" {
		errs = append(errs, response.ErrorDetail{Field: "name", Code: "required_field"})
	}
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.Phone == "" {
		errs = append(errs, response.ErrorDetail{Field: "phone", Code: "required_field"})
	}
	if req.Role == "" {
		errs = append(errs, response.ErrorDetail{Field: "role", Code: "required_field"})
	}
	if req.HireDate == "" {
		errs = append(errs, response.ErrorDetail{Field: "hire_date", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "hire_date", Code: "invalid_date"})
		return
	}

	worker, err := h.updateUC.Execute(r.Context(), workerApp.UpdateWorkerInput{
		ID:                  id,
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Role:                req.Role,
		CollectiveAgreement: req.CollectiveAgreement,
		Salary:              req.Salary,
		HireDate:            hireDate,
	})

	if err != nil {
		if err == workerDomain.ErrWorkerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "worker_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, worker)
}

func (h *WorkerHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	err := h.deleteUC.Execute(r.Context(), id)
	if err != nil {
		if err == workerDomain.ErrWorkerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "worker_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
