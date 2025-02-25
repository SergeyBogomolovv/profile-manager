package httpx

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Запрос выполнен успешно"`
}

type ErrorResponse struct {
	Code   int    `json:"code" example:"400"`
	Reason string `json:"reason" example:"Некорректный запрос"`
}

func WriteJSON(w http.ResponseWriter, payload any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, msg string, code int) error {
	return WriteJSON(w, ErrorResponse{code, msg}, code)
}

func WriteSuccess(w http.ResponseWriter, msg string, code int) error {
	return WriteJSON(w, SuccessResponse{code, msg}, code)
}
