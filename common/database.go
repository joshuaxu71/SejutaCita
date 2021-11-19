package common

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var Db *mongo.Database

func InitDb() {
	ctx := context.Background()

	var clientOptions *options.ClientOptions
	switch os.Getenv("DEPLOYMENT") {
	case "staging":
		clientOptions = options.Client().ApplyURI(os.Getenv("STAG_DB"))
	case "production":
		clientOptions = options.Client().ApplyURI(os.Getenv("PROD_DB"))
	default:
		clientOptions = options.Client().ApplyURI(os.Getenv("DEV_DB"))
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	Db = client.Database(os.Getenv("DB_NAME"))

	if collections, _ := Db.ListCollectionNames(ctx, bson.M{}); len(collections) == 0 {
		adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		generalPassword, _ := bcrypt.GenerateFromPassword([]byte("general"), bcrypt.MinCost)
		Db.Collection("users").InsertMany(
			ctx,
			[]interface{}{
				bson.M{
					"_id":         primitive.NewObjectID(),
					"created_at":  time.Now(),
					"updated_at":  time.Now(),
					"role":        "Admin",
					"first_name":  "William",
					"middle_name": "Scarra",
					"last_name":   "Lie",
					"username":    "admin",
					"password":    string(adminPassword),
				},
				bson.M{
					"_id":        primitive.NewObjectID(),
					"created_at": time.Now().AddDate(0, 0, 1),
					"updated_at": time.Now().AddDate(0, 0, 1),
					"role":       "General",
					"first_name": "Joshua",
					"username":   "general",
					"password":   string(generalPassword),
				},
			},
			nil,
		)
	}
}

func GetDb() (*mongo.Database, error) {
	err := Db.Client().Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return Db, nil
}
