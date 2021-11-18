package middleware

import (
	"SejutaCita/models"
	"context"
	"net/http"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		clientToken := r.Header.Get("token")
		if clientToken == "" {
			http.Error(rw, "No Authorization header provided", http.StatusInternalServerError)
			return
		}

		claims, err := models.ValidateToken(clientToken)
		if err != "" {
			http.Error(rw, err, http.StatusInternalServerError)
			return
		}

		var ctx context.Context
		ctx = context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "user_role", string(claims.UserRole))

		h.ServeHTTP(rw, r.WithContext(ctx))
	})
}
