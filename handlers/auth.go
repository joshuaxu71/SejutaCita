package handlers

import (
	"SejutaCita/models"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// swagger:route POST /login auth login
// Login with username and password and returns the token of the user
// responses:
//  200: tokenResponse
func Login(rw http.ResponseWriter, r *http.Request) {
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
