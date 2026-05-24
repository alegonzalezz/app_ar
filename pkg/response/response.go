package response

import (
	"encoding/json"
	"net/http"
)

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Code  string `json:"code"`
}

type WarningDetail struct {
	Field string `json:"field,omitempty"`
	Code  string `json:"code"`
}

type JSONResponse struct {
	Data     any             `json:"data"`
	Warnings []WarningDetail `json:"warnings"`
	Errors   []ErrorDetail   `json:"errors"`
}

// Success responde con un código HTTP exitoso y un payload de datos.
func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, data, nil, nil)
}

// Error responde con un código de error y los detalles correspondientes.
func Error(w http.ResponseWriter, status int, errs ...ErrorDetail) {
	JSON(w, status, nil, nil, errs)
}

// Warning responde con advertencias en la respuesta.
func Warning(w http.ResponseWriter, status int, warnings ...WarningDetail) {
	JSON(w, status, nil, warnings, nil)
}

// JSON responde con la estructura completa JSONResponse.
func JSON(w http.ResponseWriter, status int, data any, warnings []WarningDetail, errors []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	// Si los arrays vienen vacíos (nil), los dejamos tal cual o los serializamos como JSON nulo
	json.NewEncoder(w).Encode(JSONResponse{
		Data:     data,
		Warnings: warnings,
		Errors:   errors,
	})
}
