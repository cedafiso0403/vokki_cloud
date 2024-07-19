package models

import (
	"encoding/json"
	"net/http"
)

type BadRequestErrorResponse struct {
	Timestamp string `json:"timestamp" example:"2024-07-18T15:36:59Z"`
	Status    int    `json:"error_code" example:"400"`
	Message   string `json:"message" example:"Invalid request parameters"`
}

type UnauthorizedErrorResponse struct {
	Timestamp string `json:"timestamp" example:"2024-07-18T15:36:59Z"`
	Status    int    `json:"error_code" example:"401"`
	Message   string `json:"message" example:"Unauthorized access"`
}

type UserAuthenticatedResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJzdWIiOiAiMTIzNDU2Nzg5MCIsICJleHBpcmVkIjogMTY5MDY1MDAwMCwgImlhdCI6IDE2ODk3OTMwMDB9.njvE5Lgs1fjr-mL6l7QJbdFfL86D4HK4XsEFPfSb2X8"`
	TokenType string `json:"tokenType" example:"Bearer"`
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
