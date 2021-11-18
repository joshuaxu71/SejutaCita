package handlers

import (
	"SejutaCita/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUserById(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := models.UserFilter{}
	if _, ok := mux.Vars(r)["role"]; ok {
		if mux.Vars(r)["role"] == "Admin" {
			admin := models.Admin
			filter.Role = &admin
		} else {
			general := models.General
			filter.Role = &general
		}
	}
	if _, ok := mux.Vars(r)["sort"]; ok {
		if mux.Vars(r)["sort"] == "created_at.desc" {
			filter.Sort = &models.UserSort{
				Category: models.CreatedAt,
				Order:    models.Desc,
			}
		} else if mux.Vars(r)["sort"] == "first_name.asc" {
			filter.Sort = &models.UserSort{
				Category: models.FirstName,
				Order:    models.Asc,
			}
		} else if mux.Vars(r)["sort"] == "first_name.desc" {
			filter.Sort = &models.UserSort{
				Category: models.FirstName,
				Order:    models.Desc,
			}
		} else {
			filter.Sort = &models.UserSort{
				Category: models.CreatedAt,
				Order:    models.Asc,
			}
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

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	user := r.Context().Value(KeyUser{}).(models.User)
	id, err := models.CreateUser(&ctx, user)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to create user: %s", err), http.StatusInternalServerError)
	}

	rw.Write([]byte(id.Hex()))
}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	_, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "User with specified ID not found", http.StatusNotFound)
	}

	user := r.Context().Value(KeyUser{}).(models.User)

	result, err := models.UpdateUser(&ctx, mux.Vars(r)["id"], user)
	if err != nil {
		http.Error(rw, "Unable to update user", http.StatusInternalServerError)
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

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
