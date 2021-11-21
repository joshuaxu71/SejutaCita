package models

import (
	"SejutaCita/common"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Tokens that are returned in the response
// swagger:response userTokenResponse
type userTokenResponseWrapper struct {
	// in:body
	Body UserToken
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

type UserToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type SignedDetails struct {
	UserId   string
	UserRole UserRole
	jwt.StandardClaims
}

func (tokens *UserToken) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tokens)
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
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		UserId: user.Id.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}

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

	filter := bson.M{"_id": userId}

	updater := bson.M{
		"$set": bson.M{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
			"updated_at":    time.Now(),
		},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = db.Collection("users").UpdateOne(
		ctx,
		filter,
		updater,
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, _ := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt < time.Now().Unix() && claims.UserRole != "" {
		return nil, ErrExpiredToken
	}

	if claims.UserRole == "" {
		claims, err = renewTokens(claims.UserId)
		if err != nil {
			return nil, err
		}
		return claims, nil
	}

	return claims, nil
}

func renewTokens(userId string) (*SignedDetails, error) {
	ctx := context.Background()
	user, err := GetUserById(&ctx, userId)
	if err != nil {
		return nil, ErrExpiredToken
	}

	refreshToken, err := jwt.ParseWithClaims(
		*user.RefreshToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return nil, ErrExpiredToken
	}

	refreshClaims, ok := refreshToken.Claims.(*SignedDetails)
	if !ok {
		return nil, ErrExpiredToken
	}

	if refreshClaims.ExpiresAt >= time.Now().Unix() {
		signedToken, signedRefreshToken, _ := GenerateAllTokens(user)
		UpdateAllTokens(signedToken, signedRefreshToken, user.Id)
		token, err := jwt.ParseWithClaims(
			signedToken,
			&SignedDetails{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			},
		)
		if err != nil {
			return nil, ErrExpiredToken
		}
		return token.Claims.(*SignedDetails), nil
	}

	return nil, ErrExpiredToken
}
