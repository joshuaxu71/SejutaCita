package routes

import (
	"SejutaCita/handlers"
	"SejutaCita/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, l *log.Logger) {
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/user", handlers.GetUserById).
		Queries(
			"id", "{id}",
		)
	getRouter.HandleFunc("/users", handlers.GetUsers).
		Queries(
			"role", "{role}",
			"category", "{category}",
			"order", "{order}",
		)
	getRouter.Use(middleware.Middleware)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/user", handlers.CreateUser)
	postRouter.Use(handlers.MiddlewareValidateUser)
	postRouter.Use(middleware.Middleware)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/user", handlers.UpdateUser).
		Queries(
			"id", "{id}",
		)
	putRouter.Use(handlers.MiddlewareValidateUser)
	putRouter.Use(middleware.Middleware)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/user/", handlers.DeleteUser).
		Queries(
			"id", "{id}",
		)
	deleteRouter.Use(middleware.Middleware)
}
