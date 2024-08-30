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

	log.Println("Registering user")
	timeNow := time.Now().UTC()

	if r.Method != http.MethodPost {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Method not allowed",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	var newUser models.NewUserRequest

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&newUser)

	if err != nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Parameters are invalid",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	if newUser.Password == "" || newUser.ConfirmationPassword == "" || newUser.Email == "" {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Password, Email and Confirmation Password are required",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
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

	userJWT, err := utils.GenerateJWT(userCreated.ID)

	if err != nil {
		log.Println("Error generating JWT: ", err)
	} else {
		models.StoreToken(userCreated.ID, userJWT, vokki_constants.EmailToken)
		go func() {
			services.SendVerificationEmail(userCreated, userJWT)
		}()
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
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "User ID not found",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	if token == "" || token == nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Token required",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	err := services.ActivateUser(userID.(int), token.(string))

	if err != nil {
		errorResponse := httputil.UnauthorizedErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusUnauthorized,
			Message:   err.Error(),
		}
		httputil.ErrorJsonResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))

}

// Verify godoc
// @Summary Get user profile
// @Description Return profile for authenticated user
// @Tags User
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.UserProfile "User Profile"
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 401 {object} httputil.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /user [get]
func GetUser(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

	userID := r.Context().Value(vokki_constants.UserIDKey)

	if userID == 0 || userID == nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "User ID not found",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	user, err := models.GetUserProfile(userID.(int))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	httputil.SuccessJsonResponse(w, map[string]any{
		"data":       user,
		"time_stamp": utils.FormatDate(timeNow),
	})
}

// Verify godoc
// @Summary Update user profile
// @Description Update profile for authenticated user
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param User body models.UpdateUserProfileRequest false "First Name and Last Name""
// @Success 200 {object} models.UserProfile "User Profile"
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 401 {object} httputil.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /user [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

	userID := r.Context().Value(vokki_constants.UserIDKey)

	if userID == 0 || userID == nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "User ID not found",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	var updateUserProfileRequest models.UpdateUserProfileRequest

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&updateUserProfileRequest)

	if err != nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Paremeters are invalid",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	updatedUser, err := models.UpdateUserProfile(userID.(int), updateUserProfileRequest)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return

	}

	httputil.SuccessJsonResponse(w, map[string]any{
		"message":    "User updated",
		"data":       updatedUser,
		"time_stamp": utils.FormatDate(timeNow),
	})
}
