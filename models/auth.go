package models

import (
	"SejutaCita/common"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func VerifyPassword(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

type SignedDetails struct {
	UserId   string
	UserRole UserRole
	jwt.StandardClaims
}

func GenerateAllTokens(user *User) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		UserId:   user.Id.Hex(),
		UserRole: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId primitive.ObjectID) {
	db, err := common.GetDb()
	if err != nil {
		log.Panic(err)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	updateObj := bson.M{}

	updateObj["token"] = signedToken
	updateObj["refresh_token"] = signedRefreshToken

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj["updated_at"] = Updated_at

	upsert := true
	filter := bson.M{"_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = db.Collection("users").UpdateOne(
		ctx,
		filter,
		bson.M{
			"$set": updateObj,
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
