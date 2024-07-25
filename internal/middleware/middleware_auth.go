package middleware

import (
	"context"
	"net/http"
	"strings"
	"vokki_cloud/internal/auth_error"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/shared"
	"vokki_cloud/internal/utils"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if strings.TrimPrefix(authHeader, "Bearer ") == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		decodedToken, err := utils.ValidateToken(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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
			http.Error(w, auth_error.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
		ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
