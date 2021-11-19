package middleware

import (
	"SejutaCita/models"
	"context"
	"net/http"
	"strings"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		clientToken := r.Header.Get("Authorization")
		if !strings.Contains(clientToken, "Bearer") {
			http.Error(rw, "Invalid token", http.StatusInternalServerError)
			return
		}

		clientToken = strings.Replace(clientToken, "Bearer ", "", -1)

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
