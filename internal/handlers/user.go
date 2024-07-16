package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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

	var newUser models.NewUserRequest

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&newUser)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if newUser.Password == "" || newUser.ConfirmationPassword == "" || newUser.Email == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if !newUser.IsValidEmail() {
		errorResponse := models.ErrorResponse{
			Timestamp: time.Now().UTC().String(),
			Status:    http.StatusBadRequest,
			Message:   "Invalid email",
			Patch:     r.URL.Path,
		}

		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	//! Add password validation -> To define

	if newUser.Password != newUser.ConfirmationPassword {
		errorResponse := models.ErrorResponse{
			Timestamp: time.Now().UTC().String(),
			Status:    http.StatusBadRequest,
			Message:   "Passwords do not match",
			Patch:     r.URL.Path,
		}
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	user, _ := models.GetUser(newUser.Email)

	if user.Email != "" {
		errorResponse := models.ErrorResponse{
			Timestamp: time.Now().UTC().String(),
			Status:    http.StatusBadRequest,
			Message:   "Email already in use",
			Patch:     r.URL.Path,
		}
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	userCreated, err := newUser.CreateUser()

	if err != nil {
		errorResponse := models.ErrorResponse{
			Timestamp: time.Now().UTC().String(),
			Status:    http.StatusInternalServerError,
			Message:   "",
			Patch:     r.URL.Path,
		}
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

	models.SuccessJsonResponse(w, map[string]string{
		"message": "User created",
	})
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {

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
		errorResponse := models.ErrorResponse{
			Timestamp: time.Now().UTC().String(),
			Status:    http.StatusBadRequest,
			Message:   err.Error(),
			Patch:     r.URL.Path,
		}
		models.ErrorJsonResponse(w, errorResponse)
		return
	}

	w.Write([]byte(""))

}
