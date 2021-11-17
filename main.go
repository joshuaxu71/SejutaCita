package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"SejutaCita/common"
	"SejutaCita/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	common.InitDb()

	logger := log.New(os.Stdout, "log", log.LstdFlags)

	userHandler := handlers.NewUsers(logger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", userHandler.GetUsers)
	getRouter.HandleFunc("/{id:[a-z0-9]+}", userHandler.GetUserById)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", userHandler.CreateUser)
	postRouter.Use(userHandler.MiddlewareValidateProduct)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[a-z0-9]+}", userHandler.UpdateUser)
	putRouter.Use(userHandler.MiddlewareValidateProduct)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[a-z0-9]+}", userHandler.DeleteUser)

	// create a new server
	s := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      serveMux,          // set the default handler
		ErrorLog:     logger,            // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	logger.Println("Starting server on port 9090")

	err = s.ListenAndServe()
	if err != nil {
		logger.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
