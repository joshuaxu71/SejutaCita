package handlers

import (
	"SejutaCita/models"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	l *log.Logger
}

func NewAuthHandler(l *log.Logger) *AuthHandler {
	return &AuthHandler{l}
}

// swagger:route POST /login auth login
// Login with username and password and returns the token of the user
// responses:
//  200: tokenResponse
func (h *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	user := r.Context().Value(KeyUser{}).(models.User)

	existingUser, err := models.GetUserByUsername(&ctx, user.Username)
	if err != nil {
		http.Error(rw, "Incorrect credentials", http.StatusUnauthorized)
		return
	}

	if models.VerifyPassword(mux.Vars(r)["password"], user.Password) {
		http.Error(rw, "Incorrect credentials", http.StatusUnauthorized)
		return
	}

	token, refreshToken, _ := models.GenerateAllTokens(existingUser)
	models.UpdateAllTokens(token, refreshToken, existingUser.Id)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(token))
}

func (h *AuthHandler) MiddlewareValidateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		// validate the user on create
		if r.Method == http.MethodPost {
			err = user.ValidateCreate()
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error validating user: %s", err), http.StatusBadRequest)
				return
			}
		}

		// validate Role field if it's update
		if r.Method == http.MethodPut {
			err = user.ValidateUpdate()
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error validating user: %s", err), http.StatusBadRequest)
				return
			}
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
