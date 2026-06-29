package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	authApp "gcp-serverless-app/internal/auth/application"
	response "gcp-serverless-app/internal/shared/platform/http"
)

type CredentialsHandler struct {
	createAuthUC *authApp.CreateAuthUseCase
}

func NewCredentialsHandler(createAuthUC *authApp.CreateAuthUseCase) *CredentialsHandler {
	return &CredentialsHandler{createAuthUC: createAuthUC}
}

func (h *CredentialsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /administratives/{id}/credentials", h.handleCreateCredentials("administrative"))
	mux.HandleFunc("POST /workers/{id}/credentials", h.handleCreateCredentials("worker"))
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *CredentialsHandler) handleCreateCredentials(profileType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := r.PathValue("id")
		if profileID == "" {
			response.Error(w, http.StatusBadRequest, response.ErrorDetail{Field: "id", Code: "required_field"})
			return
		}

		var req credentialsRequest
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

		authID, err := generateAuthUUID()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
			return
		}

		err = h.createAuthUC.Execute(r.Context(), authApp.CreateAuthInput{
			ID:          authID,
			ProfileID:   profileID,
			ProfileType: profileType,
			Email:       req.Email,
			Password:    req.Password,
			CreatedAt:   time.Now().UTC(),
		})

		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.ErrorDetail{Code: "internal_error"})
			return
		}

		response.Success(w, http.StatusCreated, map[string]string{"message": "credenciales creadas correctamente"})
	}
}

func generateAuthUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
