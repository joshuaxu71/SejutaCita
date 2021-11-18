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
	getRouter.HandleFunc("/", handlers.GetUsers)
	getRouter.HandleFunc("/{id:[a-z0-9]+}", handlers.GetUserById)
	getRouter.Use(middleware.Middleware)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", handlers.CreateUser)
	postRouter.Use(handlers.MiddlewareValidateUser)
	postRouter.Use(middleware.Middleware)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[a-z0-9]+}", handlers.UpdateUser)
	putRouter.Use(handlers.MiddlewareValidateUser)
	putRouter.Use(middleware.Middleware)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[a-z0-9]+}", handlers.DeleteUser)
	deleteRouter.Use(middleware.Middleware)
}
