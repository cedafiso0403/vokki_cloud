package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"vokki_cloud/internal/email"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	//email.SendVerificationEmail("cedafiso@gmail.com", "1234567890")

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

	if !newUser.IsValidEmail() {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, "Invalid email", r.URL.Path)
		models.JsonResponse(w, errorResponse)
		return
	}

	//! Add password validation -> To define

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

	userCreated, err := newUser.CreateUser()

	if err != nil {
		errorResponse := models.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", r.URL.Path)
		models.JsonResponse(w, errorResponse)
		return
	}

	userJWT, err := utils.GenerateJWT(int64(userCreated.ID))

	if err != nil {
		log.Println("Error generating JWT: ", err)
	} else {
		models.StoreToken(int64(userCreated.ID), userJWT)
		email.SendVerificationEmail(userCreated, userJWT)
	}

	models.JsonResponse(w, models.NewErrorResponse(http.StatusCreated, "User created", r.URL.Path))
}
