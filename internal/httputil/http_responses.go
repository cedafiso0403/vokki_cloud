package httputil

import (
	"encoding/json"
	"net/http"
)

type UserAuthenticatedResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiAiSFMyNTeHBpcmVkIjogMTY5MDY1TMwMDB9.njvE5Lgs1fjr-mL6l7QJbdFfL86D4HK4XsEFPfSb2X8"`
	TokenType string `json:"tokenType" example:"Bearer"`
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
