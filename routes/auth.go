package routes

import (
	"SejutaCita/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router, l *log.Logger) {
	handler := handlers.NewAuthHandler(l)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/login", handler.Login)
	postRouter.Use(handler.MiddlewareValidateLogin)
}
