package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/httputil"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/services"
	"vokki_cloud/internal/utils"
)

// Register godoc
// @Summary Register an user
// @Description Register an user by email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param User body models.NewUserRequest true "Email, Password and Password Confirmation"
// @Success 200
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

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
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Invalid email",
		}

		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	//! Add password validation -> To define

	if newUser.Password != newUser.ConfirmationPassword {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Passwords do not match",
		}
		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	user, _ := models.GetUser(newUser.Email)

	if user.Email != "" {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Email already in use",
		}
		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	userCreated, err := newUser.CreateUser()

	if err != nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusInternalServerError,
			Message:   "",
		}
		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	userJWT, err := utils.GenerateJWT(int64(userCreated.ID))

	if err != nil {
		log.Println("Error generating JWT: ", err)
	} else {
		models.StoreToken(int64(userCreated.ID), userJWT, vokki_constants.EmailToken)
		services.SendVerificationEmail(userCreated, userJWT)
	}

	httputil.SuccessJsonResponse(w, map[string]string{
		"message": "User created",
	})
}

// Verify godoc
// @Summary Authenticate user
// @Description Verify user by email verification token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param        Token    query     string  true  "Email verification Token"
// @Success 200
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 401 {object} httputil.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /verify [get]
func VerifyUser(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

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

	err := services.ActivateUser(userID.(int64), token.(string))

	if err != nil {
		errorResponse := httputil.UnauthorizedErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   err.Error(),
		}
		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))

}
