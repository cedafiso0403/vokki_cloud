package middleware

import (
	"context"
	"net/http"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/httputil"
	"vokki_cloud/internal/utils"
)

func EmailVerificationMiddleware(next http.Handler) http.HandlerFunc {

	timeNow := time.Now().UTC()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				Message:   "Method not allowed",
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		//! Should be a constant?
		token := r.URL.Query().Get("token")

		if token == "" {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				Message:   "Token is required",
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		decodedToken, err := utils.ValidateToken(token)

		if err != nil {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				//Directly exposion the error message as ValidateToken returns fixed error messages set by us
				Message: err.Error(),
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
		ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
