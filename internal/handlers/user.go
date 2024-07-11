package handlers

import (
	"encoding/json"
	"net/http"
	"vokki_cloud/internal/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorResponse := models.NewErrorResponse(http.StatusMethodNotAllowed, "Method not allowed", r.URL.Path)
		models.JsonResponse(w, errorResponse)
		return
	}

	var newUser models.NewUser

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Password == "" || newUser.ConfirmationPassword == "" || newUser.Email == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	//! Add email validation

	//! Add password validation

	if newUser.Password != newUser.ConfirmationPassword {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, "Password do not match", r.URL.Path)
		models.JsonResponse(w, errorResponse)
		return
	}

	user, _ := models.GetUser(newUser.Email)

	if user.Email != "" {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, "Email already registered", r.URL.Path)
		models.JsonResponse(w, errorResponse)
		return
	}

}
