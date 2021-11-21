package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// in:body
	Body GenericError
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func (err GenericError) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(err)
}

// ErrIncorrectCredentials is an error raised when the credentials sent are incorrect
var ErrIncorrectCredentials = errors.New("incorrect credentials")

// ErrInvalidToken is an error raised when the token is invalid
var ErrInvalidToken = errors.New("invalid token")

// ErrExpiredToken is an error raised when the token has expired
var ErrExpiredToken = errors.New("expired token")

// ErrUnauthorized is an error raised when the client is not authenticated
var ErrUnauthorized = errors.New("not authenticated")

// ErrForbidden is an error raised when a user accesses operations they are not authorized for
var ErrForbidden = errors.New("insufficient access rights")

// ErrUserNotFound is an error raised when a user can not be found in the database
var ErrUserNotFound = errors.New("user not found")

// ErrDuplicateUsername is an error raised when a user with a non-unique username is being created
var ErrDuplicateUsername = errors.New("username already exists")

// ErrJsonMarshal is an error raised when server fails to marshal to json
var ErrJsonMarshal = errors.New("unable to marshal json")

// ErrJsonUnmarshal is an error raised when server fails to unmarshal from json
var ErrJsonUnmarshal = errors.New("unable to unmarshal json")
