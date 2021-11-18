package models

import (
	"SejutaCita/common"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `bson:"_id"         json:"id"`
	CreatedAt  time.Time          `bson:"created_at"  json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"  json:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at"  json:"deleted_at"`
	Role       UserRole           `bson:"role"        json:"role"         validate:"required,role"`
	FirstName  string             `bson:"first_name"  json:"first_name"   validate:"required"`
	MiddleName *string            `bson:"middle_name" json:"middle_name"`
	LastName   *string            `bson:"last_name"   json:"last_name"`
	Username   string             `bson:"username"    json:"username"     validate:"required"`
	Password   string             `bson:"password"    json:"password"     validate:"required"`
}

type Users []*User

type UserRole string

const (
	General UserRole = "General"
	Admin   UserRole = "Admin"
)

type UserSortCategory string

const (
	CreatedAt UserSortCategory = "created_at"
	FirstName UserSortCategory = "first_name"
)

type UserSort struct {
	Category UserSortCategory
	Order    SortOrder
}

type UserFilter struct {
	Role *UserRole
	Sort *UserSort
}

func (user *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("role", validateRole)

	return validate.Struct(user)
}

func validateRole(fl validator.FieldLevel) bool {
	if fl.Field().String() == string(General) || fl.Field().String() == string(Admin) {
		return true
	}
	return false
}

func (user *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(user)
}

func (users *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(users)
}

func GetUserById(ctx *context.Context, id string) (*Users, error) {
	db, err := common.GetDb()
	if err != nil {
		return nil, err
	}

	user := User{}
	filter := bson.D{primitive.E{Key: "_id", Value: common.ObjectIDFromHex(id)}}
	err = db.Collection("users").FindOne(*ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	users := append(Users{}, &user)
	return &users, nil
}

func GetUsers(ctx *context.Context, filter *UserFilter) (Users, error) {
	db, err := common.GetDb()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{}
	users := Users{}

	if filter != nil {
		if filter.Role != nil {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"role": filter.Role}})
		}
		if filter.Sort != nil {
			pipeline = append(pipeline, bson.M{"$sort": bson.M{string(filter.Sort.Category): int(filter.Sort.Order)}})
		}
	}

	cur, err := db.Collection("users").Aggregate(*ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(*ctx) {
		user := User{}
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	cur.Close(*ctx)

	return users, nil
}

func CreateUser(ctx *context.Context, user User) (primitive.ObjectID, error) {
	db, err := common.GetDb()
	if err != nil {
		return primitive.NilObjectID, err
	}

	now := time.Now()
	user.Id = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Password = common.HashAndSalt(user.Password)
	result, err := db.Collection("users").InsertOne(*ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateUser(ctx *context.Context, id string, user User) (bool, error) {
	db, err := common.GetDb()
	if err != nil {
		return false, err
	}

	_, err = GetUserById(ctx, id)
	if err != nil {
		return false, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: common.ObjectIDFromHex(id)}}
	updates := bson.M{}
	if user.Role != "" {
		updates["role"] = user.Role
	}
	if user.FirstName != "" {
		updates["first_name"] = user.FirstName
	}
	if user.MiddleName != nil {
		updates["middle_name"] = user.MiddleName
	}
	if user.LastName != nil {
		updates["last_name"] = user.LastName
	}
	if user.Password != "" {
		user.Password = common.HashAndSalt(user.Password)
		updates["password"] = user.Password
	}
	updater := bson.M{"$set": updates}

	_, err = db.Collection("users").UpdateOne(*ctx, filter, updater)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUser(ctx *context.Context, id string) (bool, error) {
	db, err := common.GetDb()
	if err != nil {
		return false, err
	}

	_, err = GetUserById(ctx, id)
	if err != nil {
		return false, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: common.ObjectIDFromHex(id)}}

	_, err = db.Collection("users").DeleteOne(*ctx, filter)
	if err != nil {
		return false, err
	}

	return true, nil
}
