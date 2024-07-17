package middleware

import (
	"context"
	"net/http"
	"time"
	"vokki_cloud/internal/auth_error"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/utils"
)

func EmailVerificationMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		//! Should be a constant?
		token := r.URL.Query().Get("token")

		if token == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		decodedToken, err := utils.ParseJWT(token)

		if err != nil {
			http.Error(w, auth_error.ErrInvalidToken.Error(), http.StatusBadRequest)
			return
		}

		if !decodedToken.VerifyExpiresAt(time.Now().Unix(), true) {
			http.Error(w, auth_error.ErrExpiredToken.Error(), http.StatusUnauthorized)
			models.RevokeToken(token)
			return
		}

		if !decodedToken.VerifyIssuer(vokki_constants.Issuer, true) {
			http.Error(w, auth_error.ErrInvalidTokenIssuer.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
		ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
