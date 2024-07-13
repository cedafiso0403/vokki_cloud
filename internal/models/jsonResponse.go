package models

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"error_code"`
	Message   string `json:"message"`
	Patch     string `json:"path"`
}

func NewErrorResponse(status int, message string, path string) ErrorResponse {
	return ErrorResponse{
		Timestamp: time.Now().UTC().String(),
		Status:    status,
		Message:   message,
		Patch:     path,
	}
}

// ErrorJsonResponse is a helper function to send a JSON response with an error message
func ErrorJsonResponse(w http.ResponseWriter, data ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(data.Status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}

//!Missing success response
