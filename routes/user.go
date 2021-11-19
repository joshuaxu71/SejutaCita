package routes

import (
	"SejutaCita/handlers"
	"SejutaCita/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, l *log.Logger) {
	handler := handlers.NewUsers(l)

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/user", handler.GetUserById).
		Queries(
			"id", "{id}",
		)
	getRouter.HandleFunc("/users", handler.GetUsers).
		Queries(
			"role", "{role}",
			"category", "{category}",
			"order", "{order}",
		)
	getRouter.Use(middleware.Middleware)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/user", handler.CreateUser)
	postRouter.Use(handler.MiddlewareValidateUser)
	postRouter.Use(middleware.Middleware)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/user", handler.UpdateUser).
		Queries(
			"id", "{id}",
		)
	putRouter.Use(handler.MiddlewareValidateUser)
	putRouter.Use(middleware.Middleware)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/user/", handler.DeleteUser).
		Queries(
			"id", "{id}",
		)
	deleteRouter.Use(middleware.Middleware)
}
