package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/email"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
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
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	//! Add password validation -> To define

	if newUser.Password != newUser.ConfirmationPassword {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, "Password do not match", r.URL.Path)
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	user, _ := models.GetUser(newUser.Email)

	if user.Email != "" {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, "Email already registered", r.URL.Path)
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	userCreated, err := newUser.CreateUser()

	if err != nil {
		errorResponse := models.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", r.URL.Path)
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	userJWT, err := utils.GenerateJWT(int64(userCreated.ID))

	if err != nil {
		log.Println("Error generating JWT: ", err)
	} else {
		models.StoreToken(int64(userCreated.ID), userJWT, vokki_constants.EmailToken)
		email.SendVerificationEmail(userCreated, userJWT)
	}

	models.ErrorJsonResponse(w, models.NewErrorResponse(http.StatusCreated, "User created", r.URL.Path))
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(vokki_constants.UserIDKey)

	token := r.Context().Value(vokki_constants.TokenKey)

	if userID == 0 || userID == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if token == "" || token == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err := models.ActivateUser(userID.(int64), token.(string))

	if err != nil {
		errorResponse := models.NewErrorResponse(http.StatusBadRequest, err.Error(), r.URL.Path)
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	w.Write([]byte("User activated"))

}
