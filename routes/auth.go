package routes

import (
	"SejutaCita/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router, l *log.Logger) {
	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/login", handlers.Login)
	postRouter.Use(handlers.MiddlewareValidateUser)
}
