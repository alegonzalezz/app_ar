package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	response "gcp-serverless-app/internal/shared/platform/http"
	userApp "gcp-serverless-app/internal/user/application"
	userDomain "gcp-serverless-app/internal/user/domain"
)

// UserHandler maneja las peticiones HTTP para usuarios.
type UserHandler struct {
	useCase *userApp.CreateUserUseCase
}

// NewUserHandler crea una nueva instancia del handler de usuarios.
func NewUserHandler(uc *userApp.CreateUserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

// CreateUserRequest define el cuerpo de la petición para crear usuario.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, response.ErrorDetail{Code: "method_not_allowed"})
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	// Validaciones básicas
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	var errs []response.ErrorDetail
	if req.Name == "" {
		errs = append(errs, response.ErrorDetail{Field: "name", Code: "required_field"})
	}
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.Password == "" {
		errs = append(errs, response.ErrorDetail{Field: "password", Code: "required_field"})
	} else if len(req.Password) < 6 {
		errs = append(errs, response.ErrorDetail{Field: "password", Code: "min_length"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	user, err := h.useCase.Execute(r.Context(), userApp.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if err == userDomain.ErrEmailAlreadyExists {
			response.Error(w, http.StatusConflict, response.ErrorDetail{Field: "email", Code: "duplicated_user"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusCreated, user)
}
