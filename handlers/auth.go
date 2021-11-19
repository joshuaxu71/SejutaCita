package handlers

import (
	"SejutaCita/models"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
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
//  200: userTokenResponse
//  401: errorResponse
//	500: errorResponse
func (h *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := ctx.Value(KeyUser{}).(models.User)

	existingUser, err := models.GetUserByUsername(&ctx, user.Username)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			rw.WriteHeader(http.StatusUnauthorized)
			models.GenericError{Message: models.ErrIncorrectCredentials.Error()}.ToJSON(rw)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		models.GenericError{Message: err.Error()}.ToJSON(rw)
		return
	}

	if models.VerifyPassword(mux.Vars(r)["password"], user.Password) {
		rw.WriteHeader(http.StatusUnauthorized)
		models.GenericError{Message: models.ErrIncorrectCredentials.Error()}.ToJSON(rw)
		return
	}

	token, refreshToken, _ := models.GenerateAllTokens(existingUser)
	models.UpdateAllTokens(token, refreshToken, existingUser.Id)

	rw.WriteHeader(http.StatusOK)
	tokens := models.UserToken{
		Token:        token,
		RefreshToken: refreshToken,
	}
	tokens.ToJSON(rw)
}

func (h *AuthHandler) MiddlewareValidateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			models.GenericError{Message: models.ErrJsonUnmarshal.Error()}.ToJSON(rw)
			return
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
