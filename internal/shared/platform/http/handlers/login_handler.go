package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	authApp "gcp-serverless-app/internal/auth/application"
	authDomain "gcp-serverless-app/internal/auth/domain"
	response "gcp-serverless-app/internal/shared/platform/http"
)

// LoginHandler maneja las peticiones HTTP de login.
type LoginHandler struct {
	useCase *authApp.LoginUseCase
}

// NewLoginHandler crea una nueva instancia del handler de login.
func NewLoginHandler(uc *authApp.LoginUseCase) *LoginHandler {
	return &LoginHandler{useCase: uc}
}

// LoginRequest define el cuerpo de la petición de login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, response.ErrorDetail{Code: "method_not_allowed"})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	var errs []response.ErrorDetail
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.Password == "" {
		errs = append(errs, response.ErrorDetail{Field: "password", Code: "required_field"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	userInfo, err := h.useCase.Execute(r.Context(), authApp.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if err == authDomain.ErrUserNotFound {
			response.Error(w, http.StatusUnauthorized, response.ErrorDetail{Code: "invalid_credentials"})
			return
		}
		if err == authDomain.ErrInvalidCredentials {
			response.Error(w, http.StatusUnauthorized, response.ErrorDetail{Code: "invalid_credentials"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, userInfo)
}
