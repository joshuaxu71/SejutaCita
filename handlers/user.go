package handlers

import (
	"SejutaCita/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route GET /user user getUserById
// Returns a user by ID
// responses:
//  200: userResponse
func GetUserById(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	err = user.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /users users getUsers
// Returns all users with optional filter and sorting
// responses:
//  200: usersResponse
func GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		http.Error(rw, "Unable to retrieve data", http.StatusInternalServerError)
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route POST /user user createUser
// Inserts a User in the database and returns the ID of the created User
// responses:
//  200: userIdResponse
func CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	user := r.Context().Value(KeyUser{}).(models.UserCreate)
	id, err := models.CreateUser(&ctx, user)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to create user: %s", err), http.StatusInternalServerError)
	}

	rw.Write([]byte(id.Hex()))
}

// swagger:route PUT /user user updateUser
// Updates a User in the database and returns a boolean based on the success of the update
// responses:
//  200: booleanResponse
func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	_, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "User with specified ID not found", http.StatusNotFound)
	}

	user := r.Context().Value(KeyUser{}).(models.UserUpdate)

	result, err := models.UpdateUser(&ctx, mux.Vars(r)["id"], user)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

// swagger:route DELETE /user user deleteUser
// Deletes a User in the database and returns a boolean based on the success of the update
// responses:
//  200: booleanResponse
func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	result, err := models.DeleteUser(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

type KeyUser struct{}

func MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		// validate the user if it's create
		if r.Method == http.MethodPost && r.RequestURI != "/login" {
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
