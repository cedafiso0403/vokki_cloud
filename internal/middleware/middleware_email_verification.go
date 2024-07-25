package middleware

import (
	"context"
	"net/http"
	vokki_constants "vokki_cloud/internal/constants"
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

		decodedToken, err := utils.ValidateToken(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, vokki_constants.UserIDKey, decodedToken.UserID)
		ctx = context.WithValue(ctx, vokki_constants.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
