package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"SejutaCita/common"
	"SejutaCita/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	common.InitDb()

	l := log.New(os.Stdout, "log", log.LstdFlags)

	r := mux.NewRouter()
	routes.AuthRoutes(r, l)
	routes.UserRoutes(r, l)

	// create a new server
	s := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	l.Println("Starting server on port 9090")

	err = s.ListenAndServe()
	if err != nil {
		l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
