package handlers

import (
	"SejutaCita/models"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

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

func CreateToken(userId string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
