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
			rw.WriteHeader(http.StatusUnauthorized)
			models.GenericError{Message: models.ErrUnauthorized.Error()}.ToJSON(rw)
			return
		}

		clientToken = strings.Replace(clientToken, "Bearer ", "", -1)

		claims, err := models.ValidateToken(clientToken)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			models.GenericError{Message: err.Error()}.ToJSON(rw)
			return
		}

		var ctx context.Context
		ctx = context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "user_role", string(claims.UserRole))

		h.ServeHTTP(rw, r.WithContext(ctx))
	})
}
