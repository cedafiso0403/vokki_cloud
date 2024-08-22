package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/httputil"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/shared"
	"vokki_cloud/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {

	timeNow := time.Now().UTC()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				Message:   "Authorization header is required",
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		if strings.TrimPrefix(authHeader, "Bearer ") == "" {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				Message:   "Token is required",
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		decodedToken, err := utils.ValidateToken(token)

		if err != nil {

			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				//Directly exposion the error message as ValidateToken returns fixed error message set by us
				Message: err.Error(),
			}
			httputil.ErrorJsonResponse(w, errorResponse, errorResponse.Status)
			return
		}

		if shared.GetTokenManager().TokenExists(token) {

			ctx := r.Context()
			ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
			ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		if !models.VerifyToken(token) {
			errorResponse := httputil.UnauthorizedErrorResponse{
				Timestamp: utils.FormatDate(timeNow),
				Status:    http.StatusUnauthorized,
				Message:   "Token is invalid",
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
