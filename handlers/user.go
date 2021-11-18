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

type Users struct {
	l *log.Logger
}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) GetUserById(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User", mux.Vars(r)["id"])
	ctx := context.TODO()

	users, err := models.GetUserById(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "User with specified ID not found", http.StatusNotFound)
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")
	ctx := context.TODO()

	users, err := models.GetUsers(&ctx, nil)
	if err != nil {
		http.Error(rw, "Unable to retrieve data", http.StatusInternalServerError)
	}

	err = users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (u *Users) CreateUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")
	ctx := context.TODO()

	user := r.Context().Value(KeyUser{}).(models.User)
	id, err := models.CreateUser(&ctx, user)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to create user: %s", err), http.StatusInternalServerError)
	}

	rw.Write([]byte(id.Hex()))
}

func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle PUT User", mux.Vars(r)["id"])
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

func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle DELETE User", mux.Vars(r)["id"])
	ctx := context.TODO()

	result, err := models.DeleteUser(&ctx, mux.Vars(r)["id"])
	if err != nil {
		http.Error(rw, "User with specified ID not found", http.StatusNotFound)
	}

	rw.Write([]byte(strconv.FormatBool(result)))
}

type KeyUser struct{}

func (u Users) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		// validate the user if it's create
		if r.Method == http.MethodPost {
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
