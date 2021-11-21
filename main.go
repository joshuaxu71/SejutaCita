// Package classification of SejutaCita
//
// Documentation for SejutaCita
//
//	Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//  SecurityDefinitions:
//  api_key:
//   description: "JWT Authorization header using the Bearer scheme."
//   type: apiKey
//   name: Authorization
//   in: header
// swagger:meta
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"SejutaCita/common"
	"SejutaCita/routes"
)

func main() {
	// loading the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initializing the database
	common.InitDb()

	// create the logger
	l := log.New(os.Stdout, "SejutaCita: ", log.LstdFlags)

	// create the router and serve the swagger documentation
	r := mux.NewRouter()
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Methods(http.MethodGet).Subrouter().Handle("/docs", sh)
	r.Methods(http.MethodGet).Subrouter().Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// add routes to the router
	routes.AuthRoutes(r, l)
	routes.UserRoutes(r, l)

	// create a new server
	s := http.Server{
		Addr:         os.Getenv("PORT"), // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	l.Printf("Starting server on port %s", os.Getenv("PORT"))

	err = s.ListenAndServe()
	if err != nil {
		l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
