package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	adminApp "gcp-serverless-app/internal/administrative/application"
	adminDomain "gcp-serverless-app/internal/administrative/domain"
	response "gcp-serverless-app/internal/shared/platform/http"
)

type AdministrativeHandler struct {
	createUC *adminApp.CreateAdministrativeUseCase
	getUC    *adminApp.GetAdministrativeUseCase
	listUC   *adminApp.ListAdministrativesUseCase
	updateUC *adminApp.UpdateAdministrativeUseCase
	deleteUC *adminApp.DeleteAdministrativeUseCase
}

func NewAdministrativeHandler(
	createUC *adminApp.CreateAdministrativeUseCase,
	getUC *adminApp.GetAdministrativeUseCase,
	listUC *adminApp.ListAdministrativesUseCase,
	updateUC *adminApp.UpdateAdministrativeUseCase,
	deleteUC *adminApp.DeleteAdministrativeUseCase,
) *AdministrativeHandler {
	return &AdministrativeHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *AdministrativeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /administratives", h.handleCreate)
	mux.HandleFunc("GET /administratives", h.handleList)
	mux.HandleFunc("GET /administratives/{id}", h.handleGet)
	mux.HandleFunc("PUT /administratives/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /administratives/{id}", h.handleDelete)
}

type administrativeRequest struct {
	Name                string   `json:"name"`
	Email               string   `json:"email"`
	Phone               string   `json:"phone"`
	Role                string   `json:"role"`
	CollectiveAgreement *string  `json:"collective_agreement"`
	WorkSchedule        string   `json:"work_schedule"`
	Salary              *float64 `json:"salary"`
	HireDate            string   `json:"hire_date"`
}

func (h *AdministrativeHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req administrativeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Role = strings.TrimSpace(req.Role)
	req.WorkSchedule = strings.TrimSpace(req.WorkSchedule)
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
	if req.WorkSchedule == "" {
		errs = append(errs, response.ErrorDetail{Field: "work_schedule", Code: "required_field"})
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

	admin, err := h.createUC.Execute(r.Context(), adminApp.CreateAdministrativeInput{
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Role:                req.Role,
		CollectiveAgreement: req.CollectiveAgreement,
		WorkSchedule:        req.WorkSchedule,
		Salary:              req.Salary,
		HireDate:            hireDate,
	})

	if err != nil {
		if err == adminDomain.ErrAdministrativeEmailExists {
			response.Error(w, http.StatusConflict, response.ErrorDetail{Field: "email", Code: "already_exists"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, admin)
}

func (h *AdministrativeHandler) handleList(w http.ResponseWriter, r *http.Request) {
	administratives, err := h.listUC.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	if administratives == nil {
		administratives = []*adminDomain.Administrative{}
	}

	response.Success(w, http.StatusOK, administratives)
}

func (h *AdministrativeHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	admin, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == adminDomain.ErrAdministrativeNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "administrative_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, admin)
}

func (h *AdministrativeHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req administrativeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Role = strings.TrimSpace(req.Role)
	req.WorkSchedule = strings.TrimSpace(req.WorkSchedule)
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
	if req.WorkSchedule == "" {
		errs = append(errs, response.ErrorDetail{Field: "work_schedule", Code: "required_field"})
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

	admin, err := h.updateUC.Execute(r.Context(), adminApp.UpdateAdministrativeInput{
		ID:                  id,
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Role:                req.Role,
		CollectiveAgreement: req.CollectiveAgreement,
		WorkSchedule:        req.WorkSchedule,
		Salary:              req.Salary,
		HireDate:            hireDate,
	})

	if err != nil {
		if err == adminDomain.ErrAdministrativeNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "administrative_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, admin)
}

func (h *AdministrativeHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	err := h.deleteUC.Execute(r.Context(), id)
	if err != nil {
		if err == adminDomain.ErrAdministrativeNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "administrative_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
