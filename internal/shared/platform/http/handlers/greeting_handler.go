package handlers

import (
	"net/http"

	greetApp "gcp-serverless-app/internal/greeting/application"
	response "gcp-serverless-app/internal/shared/platform/http"
)

// GreetingHandler maneja las peticiones HTTP de saludos.
type GreetingHandler struct {
	useCase *greetApp.FindGreetingUseCase
}

// NewGreetingHandler crea una nueva instancia del handler de saludos.
func NewGreetingHandler(uc *greetApp.FindGreetingUseCase) *GreetingHandler {
	return &GreetingHandler{useCase: uc}
}

func (h *GreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	result, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.ErrorDetail{Code: "not_found"})
		return
	}
	response.Success(w, http.StatusOK, result)
}
