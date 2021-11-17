package common

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Db *mongo.Client

func InitDb() {
	var clientOptions *options.ClientOptions

	switch os.Getenv("deployment") {
	case "staging":
		clientOptions = options.Client().ApplyURI(os.Getenv("stag_db"))
	case "production":
		clientOptions = options.Client().ApplyURI(os.Getenv("prod_db"))
	default:
		clientOptions = options.Client().ApplyURI(os.Getenv("dev_db"))
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	Db = client
}

func GetDb() *mongo.Client {
	err := Db.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}

	return Db
}
