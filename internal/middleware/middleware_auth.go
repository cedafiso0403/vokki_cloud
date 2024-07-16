package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
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

		decodedToken, err := utils.ParseJWT(token)

		if err != nil {
			http.Error(w, "Token not valid", http.StatusBadRequest)
			return
		}

		if !decodedToken.VerifyExpiresAt(time.Now().Unix(), true) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			models.RevokeToken(token)
			return
		}

		if !decodedToken.VerifyIssuer(vokki_constants.Issuer, true) {
			http.Error(w, "Invalid token issuer", http.StatusUnauthorized)
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
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
		ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
