package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	authApp "gcp-serverless-app/internal/auth/application"
	authDomain "gcp-serverless-app/internal/auth/domain"
	response "gcp-serverless-app/internal/shared/platform/http"
)

// ChangePasswordHandler maneja las peticiones HTTP de cambio de contraseña.
type ChangePasswordHandler struct {
	useCase *authApp.ChangePasswordUseCase
}

// NewChangePasswordHandler crea una nueva instancia del handler de cambio de contraseña.
func NewChangePasswordHandler(uc *authApp.ChangePasswordUseCase) *ChangePasswordHandler {
	return &ChangePasswordHandler{useCase: uc}
}

// ChangePasswordRequest define el cuerpo de la petición de cambio de contraseña.
type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, response.ErrorDetail{Code: "method_not_allowed"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrorDetail{Code: "invalid_json"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	var errs []response.ErrorDetail
	if req.Email == "" {
		errs = append(errs, response.ErrorDetail{Field: "email", Code: "required_field"})
	}
	if req.OldPassword == "" {
		errs = append(errs, response.ErrorDetail{Field: "old_password", Code: "required_field"})
	}
	if req.NewPassword == "" {
		errs = append(errs, response.ErrorDetail{Field: "new_password", Code: "required_field"})
	} else if len(req.NewPassword) < 6 {
		errs = append(errs, response.ErrorDetail{Field: "new_password", Code: "min_length"})
	}

	if len(errs) > 0 {
		response.Error(w, http.StatusBadRequest, errs...)
		return
	}

	err := h.useCase.Execute(r.Context(), authApp.ChangePasswordInput{
		Email:       req.Email,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		if err == authDomain.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, response.ErrorDetail{Field: "email", Code: "user_not_found"})
			return
		}
		if err == authDomain.ErrInvalidCredentials {
			response.Error(w, http.StatusUnauthorized, response.ErrorDetail{Field: "old_password", Code: "invalid_user"})
			return
		}
		response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
		return
	}

	response.Success(w, http.StatusOK, "contraseña actualizada exitosamente")
}
