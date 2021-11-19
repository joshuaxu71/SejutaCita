package models

import (
	"SejutaCita/common"
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// A token that is returned in the response
// swagger:response tokenResponse
type tokenResponseWrapper struct {
	// in:body
	Body struct {
		Token string
	}
}

// swagger:parameters login
type loginParameterWrapper struct {
	// The username and password of the user
	// in:body
	Body struct {
		// required:true
		Username string
		// required:true
		Password string
	}
}

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

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
