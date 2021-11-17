package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringAddress(s string) *string {
	return &s
}

func ObjectIDFromHex(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}
