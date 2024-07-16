package models

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"error_code"`
	Message   string `json:"message"`
	Patch     string `json:"path"`
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

func SuccessJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}
