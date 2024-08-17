package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/httputil"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/services"
	"vokki_cloud/internal/utils"
)

// Login godoc
// @Summary Authenticate user
// @Description Authenticate user by email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param Credentials body services.Credentials true "Email and Password"
// @Success 200 {object} httputil.UserAuthenticatedResponse "Success"
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 401 {object} httputil.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

	credentials := services.Credentials{}

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&credentials)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Password and Email are required",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	userID, token, err := services.Authenticate(credentials)

	if err != nil {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Incorrect Password or Email",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	err = models.StoreToken(int64(userID), token, vokki_constants.AuthToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.SuccessJsonResponse(w, httputil.UserAuthenticatedResponse{Token: token, TokenType: "Bearer"})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Send email with token to reset password if user exists, is active and is not using provider authentication, otherwise does nothing
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body services.NewPasswordEmailRequest true "User Email"
// @Success 200
// @Failure 400 {object} httputil.BadRequestErrorResponse "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /reset-password [post]
func RequestResetPassword(w http.ResponseWriter, r *http.Request) {

	timeNow := time.Now().UTC()

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var userEmail services.NewPasswordEmailRequest

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&userEmail)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if userEmail.Email == "" {
		errorResponse := httputil.BadRequestErrorResponse{
			Timestamp: utils.FormatDate(timeNow),
			Status:    http.StatusBadRequest,
			Message:   "Email is required",
		}
		httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
		return
	}

	err = services.RequestNewPasswordEmail(userEmail.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))

}
