package httputil

import (
	"encoding/json"
	"net/http"
)

type BadRequestErrorResponse struct {
	Status    int    `json:"error_code" example:"400"`
	Message   string `json:"message" example:"Invalid request parameters"`
	Timestamp string `json:"timestamp" example:"2024-07-18T15:36:59Z"`
}

type UnauthorizedErrorResponse struct {
	Status    int    `json:"error_code" example:"401"`
	Message   string `json:"message" example:"Unauthorized access"`
	Timestamp string `json:"timestamp" example:"2024-07-18T15:36:59Z"`
}

func ErrorJsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}
