package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	customerApp "gcp-serverless-app/internal/customer/application"
	customerDomain "gcp-serverless-app/internal/customer/domain"
	response "gcp-serverless-app/internal/shared/platform/http"
)

type CustomerHandler struct {
	createUC *customerApp.CreateCustomerUseCase
	getUC    *customerApp.GetCustomerUseCase
	listUC   *customerApp.ListCustomersUseCase
	updateUC *customerApp.UpdateCustomerUseCase
	deleteUC *customerApp.DeleteCustomerUseCase
}

func NewCustomerHandler(
	createUC *customerApp.CreateCustomerUseCase,
	getUC *customerApp.GetCustomerUseCase,
	listUC *customerApp.ListCustomersUseCase,
	updateUC *customerApp.UpdateCustomerUseCase,
	deleteUC *customerApp.DeleteCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

// RegisterRoutes registra las rutas de clientes usando ServeMux.
func (h *CustomerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", h.handleCreate)
	mux.HandleFunc("GET /customers", h.handleList)
	mux.HandleFunc("GET /customers/{id}", h.handleGet)
	mux.HandleFunc("PUT /customers/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /customers/{id}", h.handleDelete)
}

type customerRequest struct {
	Name             string  `json:"name"`
	PhoneNumber      string  `json:"phone_number"`
	ExtraPhoneNumber *string `json:"extra_phone_number"`
	ContactEmail     string  `json:"contact_email"`
	ManagerName      string  `json:"manager_name"`
	Address          string  `json:"address"`
}

func (h *CustomerHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	req.ContactEmail = strings.TrimSpace(req.ContactEmail)
	req.ManagerName = strings.TrimSpace(req.ManagerName)
	req.Address = strings.TrimSpace(req.Address)

	var errs []response.ErrorDetail
	if req.Name == "" {
		errs = append(errs, response.ErrorDetail{Field: "name", Code: "required_field"})
	}
	if req.PhoneNumber == "" {
		errs = append(errs, response.ErrorDetail{Field: "phone_number", Code: "required_field"})
	}
	if req.ContactEmail == "" {
		errs = append(errs, response.ErrorDetail{Field: "contact_email", Code: "required_field"})
	}
	if req.ManagerName == "" {
		errs = append(errs, response.ErrorDetail{Field: "manager_name", Code: "required_field"})
	}
	if req.Address == "" {
		errs = append(errs, response.ErrorDetail{Field: "address", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	customer, err := h.createUC.Execute(r.Context(), customerApp.CreateCustomerInput{
		Name:             req.Name,
		PhoneNumber:      req.PhoneNumber,
		ExtraPhoneNumber: req.ExtraPhoneNumber,
		ContactEmail:     req.ContactEmail,
		ManagerName:      req.ManagerName,
		Address:          req.Address,
	})

	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) handleList(w http.ResponseWriter, r *http.Request) {
	customers, err := h.listUC.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	// Si es nulo, retornar slice vacío en lugar de null
	if customers == nil {
		customers = []*customerDomain.Customer{}
	}

	response.Success(w, http.StatusOK, customers)
}

func (h *CustomerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	customer, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		if err == customerDomain.ErrCustomerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "customer_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, customer)
}

func (h *CustomerHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	req.ContactEmail = strings.TrimSpace(req.ContactEmail)
	req.ManagerName = strings.TrimSpace(req.ManagerName)
	req.Address = strings.TrimSpace(req.Address)

	var errs []response.ErrorDetail
	if req.Name == "" {
		errs = append(errs, response.ErrorDetail{Field: "name", Code: "required_field"})
	}
	if req.PhoneNumber == "" {
		errs = append(errs, response.ErrorDetail{Field: "phone_number", Code: "required_field"})
	}
	if req.ContactEmail == "" {
		errs = append(errs, response.ErrorDetail{Field: "contact_email", Code: "required_field"})
	}
	if req.ManagerName == "" {
		errs = append(errs, response.ErrorDetail{Field: "manager_name", Code: "required_field"})
	}
	if req.Address == "" {
		errs = append(errs, response.ErrorDetail{Field: "address", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	customer, err := h.updateUC.Execute(r.Context(), customerApp.UpdateCustomerInput{
		ID:               id,
		Name:             req.Name,
		PhoneNumber:      req.PhoneNumber,
		ExtraPhoneNumber: req.ExtraPhoneNumber,
		ContactEmail:     req.ContactEmail,
		ManagerName:      req.ManagerName,
		Address:          req.Address,
	})

	if err != nil {
		if err == customerDomain.ErrCustomerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "customer_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, customer)
}

func (h *CustomerHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
		return
	}

	err := h.deleteUC.Execute(r.Context(), id)
	if err != nil {
		if err == customerDomain.ErrCustomerNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "customer_not_found"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
