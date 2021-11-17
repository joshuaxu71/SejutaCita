package main

import (
	"log"

	"github.com/joho/godotenv"

	"SejutaCita/common"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	common.InitDb()
}
