package handlers

import (
	"SejutaCita/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	l *log.Logger
}

func NewUserHandler(l *log.Logger) *UserHandler {
	return &UserHandler{l}
}

// swagger:route GET /user user getUserById
// Returns a user by ID
// responses:
//  200: userResponse
//  401: errorResponse
//  403: errorResponse
//  404: errorResponse
//  500: errorResponse
func (h *UserHandler) GetUserById(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("user_role") != "Admin" {
		if ctx.Value("user_id") != mux.Vars(r)["id"] {
			rw.WriteHeader(http.StatusForbidden)
			models.GenericError{Message: models.ErrForbidden.Error()}.ToJSON(rw)
			return
		}
	}

	user, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			rw.WriteHeader(http.StatusNotFound)
			models.GenericError{Message: models.ErrUserNotFound.Error()}.ToJSON(rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			models.GenericError{Message: fmt.Sprintf("Unable to get user: %s", err)}.ToJSON(rw)
			return
		}
	}

	err = user.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		models.GenericError{Message: models.ErrJsonMarshal.Error()}.ToJSON(rw)
		return
	}
}

// swagger:route GET /users users getUsers
// Returns all users with optional filter and sorting
// responses:
//  200: usersResponse
//  401: errorResponse
//	403: errorResponse
//  404: errorResponse
//  500: errorResponse
func (h *UserHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("user_role") != "Admin" {
		rw.WriteHeader(http.StatusForbidden)
		models.GenericError{Message: models.ErrForbidden.Error()}.ToJSON(rw)
		return
	}

	filter := models.UserFilter{
		Sort: &models.UserSort{},
	}
	if _, ok := mux.Vars(r)["role"]; ok {
		if mux.Vars(r)["role"] == "Admin" {
			admin := models.Admin
			filter.Role = &admin
		} else {
			general := models.General
			filter.Role = &general
		}
	}
	if _, ok := mux.Vars(r)["category"]; ok {
		if mux.Vars(r)["category"] == "first_name" {
			filter.Sort.Category = models.FirstName
		} else {
			filter.Sort.Category = models.CreatedAt
		}
	}
	if _, ok := mux.Vars(r)["order"]; ok {
		if mux.Vars(r)["order"] == strconv.Itoa(int(models.Desc)) {
			filter.Sort.Order = models.Desc - 1
		} else {
			filter.Sort.Order = models.Asc - 1
		}
	}

	users, err := models.GetUsers(&ctx, &filter)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			rw.WriteHeader(http.StatusNotFound)
			models.GenericError{Message: models.ErrUserNotFound.Error()}.ToJSON(rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			models.GenericError{Message: fmt.Sprintf("Unable to get users: %s", err)}.ToJSON(rw)
			return
		}
	}

	err = users.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		models.GenericError{Message: models.ErrJsonMarshal.Error()}.ToJSON(rw)
		return
	}
}

// swagger:route POST /user user createUser
// Inserts a User in the database and returns the ID of the created User
// responses:
//  200: userIdResponse
//  401: errorResponse
//	403: errorResponse
//	409: errorResponse
//  500: errorResponse
func (h *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("user_role") != "Admin" {
		rw.WriteHeader(http.StatusForbidden)
		models.GenericError{Message: models.ErrForbidden.Error()}.ToJSON(rw)
		return
	}

	user := r.Context().Value(KeyUser{}).(models.User)
	id, err := models.CreateUser(&ctx, user)
	if err != nil {
		switch err {
		case models.ErrDuplicateUsername:
			rw.WriteHeader(http.StatusConflict)
			models.GenericError{Message: err.Error()}.ToJSON(rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			models.GenericError{Message: fmt.Sprintf("Unable to create user: %s", err)}.ToJSON(rw)
			return
		}
	}

	rw.Write([]byte(id.Hex()))
}

// swagger:route PUT /user user updateUser
// Updates a User in the database and returns a boolean based on the success of the update
// responses:
//  200: booleanResponse
//  401: errorResponse
//	403: errorResponse
//  404: errorResponse
//  500: errorResponse
func (u *UserHandler) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("user_role") != "Admin" {
		rw.WriteHeader(http.StatusForbidden)
		models.GenericError{Message: models.ErrForbidden.Error()}.ToJSON(rw)
		return
	}

	user := r.Context().Value(KeyUser{}).(models.User)

	result, err := models.UpdateUser(&ctx, mux.Vars(r)["id"], user)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			rw.WriteHeader(http.StatusNotFound)
			models.GenericError{Message: models.ErrUserNotFound.Error()}.ToJSON(rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			models.GenericError{Message: fmt.Sprintf("Unable to update user: %s", err)}.ToJSON(rw)
			return
		}
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

// swagger:route DELETE /user user deleteUser
// Deletes a User in the database and returns a boolean based on the success of the update
// responses:
//  200: booleanResponse
//  401: errorResponse
//	403: errorResponse
//  404: errorResponse
//  500: errorResponse
func (h *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("user_role") != "Admin" {
		rw.WriteHeader(http.StatusForbidden)
		models.GenericError{Message: models.ErrForbidden.Error()}.ToJSON(rw)
		return
	}

	result, err := models.DeleteUser(&ctx, mux.Vars(r)["id"])
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			rw.WriteHeader(http.StatusNotFound)
			models.GenericError{Message: models.ErrUserNotFound.Error()}.ToJSON(rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			models.GenericError{Message: fmt.Sprintf("Unable to delete user: %s", err)}.ToJSON(rw)
			return
		}
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

type KeyUser struct{}

func (h *UserHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			models.GenericError{Message: models.ErrJsonUnmarshal.Error()}.ToJSON(rw)
			return
		}

		// validate the user on create
		if r.Method == http.MethodPost {
			err = user.ValidateCreate()
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				models.GenericError{Message: fmt.Sprintf("Error validating user: %s", err)}.ToJSON(rw)
				return
			}
		}

		// validate Role field if it's update
		if r.Method == http.MethodPut {
			err = user.ValidateUpdate()
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				models.GenericError{Message: fmt.Sprintf("Error validating user: %s", err)}.ToJSON(rw)
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
