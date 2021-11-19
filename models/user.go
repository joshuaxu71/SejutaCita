package models

import (
	"SejutaCita/common"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A user that is returned in the response
// swagger:response userResponse
type userResponseWrapper struct {
	// in:body
	Body User
}

// Users that are returned in the response
// swagger:response usersResponse
type usersResponseWrapper struct {
	// in:body
	Body []User
}

// User ID (string) that is returned in the response
// swagger:response userIdResponse
type userIdResponseWrapper struct {
	// in:body
	Body struct {
		Id primitive.ObjectID
	}
}

// swagger:parameters getUserById updateUser deleteUser
type userIdParameterWrapper struct {
	// The ID of the user to perform the operation on
	// in:query
	// required:true
	Id string `json:"id"`
}

// swagger:parameters getUsers
type usersGetParameterWrapper struct {
	// The filter based on user's role
	// in:query
	Role UserRole `json:"role"`
	// The sorting based on category
	// in:query
	Category UserSortCategory `json:"category"`
	// The sorting based on order
	// in:query
	Order SortOrder `json:"order"`
}

// swagger:parameters createUser
type userCreateParameterWrapper struct {
	// The details of the User that will be created
	// in:body
	// required:true
	Body UserCreate
}

// swagger:parameters updateUser
type userUpdateParameterWrapper struct {
	// The details of the User that will be updated
	// in:body
	// required:true
	Body UserUpdate
}

// User defines the structure for an API User on GET methods
// swagger:model
type User struct {
	// the ID of the user
	// required:true
	// swagger:strfmt bsonobjectid
	Id primitive.ObjectID `bson:"_id"           json:"id"`
	// the date the user was created at
	// required:true
	CreatedAt time.Time `bson:"created_at"    json:"created_at"`
	// the date the user was last updated at
	// required:true
	UpdatedAt time.Time `bson:"updated_at"    json:"updated_at"`
	// the date the user was deleted at
	DeletedAt *time.Time `bson:"deleted_at"    json:"deleted_at"`
	// the token of the user
	Token *string `bson:"token"         json:"token"`
	// the refresh token of the user
	RefreshToken *string `bson:"refresh_token" json:"refresh_token"`
	// the role of the user
	// required:true
	Role UserRole `bson:"role"          json:"role"         validate:"role"`
	// the first name of the user
	// required:true
	FirstName string `bson:"first_name"    json:"first_name"   validate:"first_name"`
	// the middle name of the user
	MiddleName *string `bson:"middle_name"   json:"middle_name"`
	// the last name of the user
	LastName *string `bson:"last_name"     json:"last_name"`
	// the username of the user
	// required:true
	Username string `bson:"username"      json:"username"     validate:"username"`
	// the password of the user
	// required:true
	Password string `bson:"password"      json:"password"     validate:"password"`
}

// UserCreate defines the structure for an API User on POST methods
// swagger:model
type UserCreate struct {
	// the role of the user
	// required:true
	Role UserRole `bson:"role"          json:"role"         validate:"role"`
	// the first name of the user
	// required:true
	FirstName string `bson:"first_name"    json:"first_name"   validate:"first_name"`
	// the middle name of the user
	MiddleName *string `bson:"middle_name"   json:"middle_name"`
	// the last name of the user
	LastName *string `bson:"last_name"     json:"last_name"`
	// the username of the user
	// required:true
	Username string `bson:"username"      json:"username"     validate:"username"`
	// the password of the user
	// required:true
	Password string `bson:"password"      json:"password"     validate:"password"`
}

// UserUpdate defines the structure for an API User on PUT methods
// swagger:model
type UserUpdate struct {
	// the role of the user
	Role UserRole `bson:"role"          json:"role"         validate:"role"`
	// the first name of the user
	FirstName string `bson:"first_name"    json:"first_name"   validate:"first_name"`
	// the middle name of the user
	MiddleName *string `bson:"middle_name"   json:"middle_name"`
	// the last name of the user
	LastName *string `bson:"last_name"     json:"last_name"`
	// the password of the user
	Password string `bson:"password"      json:"password"     validate:"password"`
}

type Users []*User

// swagger:enum UserRole
type UserRole string

const (
	General UserRole = "General"
	Admin   UserRole = "Admin"
)

// swagger:enum UserSortCategory
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

func (user *User) ValidateCreate() error {
	validate := validator.New()
	validate.RegisterValidation("role", validateRole)
	validate.RegisterValidation("first_name", validateFirstName)
	validate.RegisterValidation("username", validateUsername)
	validate.RegisterValidation("password", validatePassword)

	return validate.Struct(user)
}

func (user *User) ValidateUpdate() error {
	validate := validator.New()
	validate.RegisterValidation("role", validateRole)
	validate.RegisterValidation("first_name", EmptyValidate)
	validate.RegisterValidation("username", EmptyValidate)
	validate.RegisterValidation("password", EmptyValidate)

	return validate.Struct(user)
}

func validateRole(fl validator.FieldLevel) bool {
	return fl.Field().String() == string(General) || fl.Field().String() == string(Admin)
}

func validateFirstName(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

func validateUsername(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

func validatePassword(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

func (user *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(user)
}

func (user *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(user)
}

func (users *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(users)
}

func GetUserById(ctx *context.Context, id string) (*User, error) {
	db, err := common.GetDb()
	if err != nil {
		return nil, err
	}
	user := User{}
	filter := bson.M{"_id": common.ObjectIDFromHex(id)}
	err = db.Collection("users").FindOne(*ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(ctx *context.Context, username string) (*User, error) {
	db, err := common.GetDb()
	if err != nil {
		return nil, err
	}

	user := User{}
	filter := bson.M{"username": username}
	err = db.Collection("users").FindOne(*ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func GetUsers(ctx *context.Context, filter *UserFilter) (Users, error) {
	context := *ctx

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

	cur, err := db.Collection("users").Aggregate(*ctx, pipeline, options.Aggregate().SetCollation(&options.Collation{Locale: "en"}))
	if err != nil {
		if err == mongo.ErrNilCursor {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	for cur.Next(*ctx) {
		user := User{}
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		if context.Value("user_role") != "Admin" {
			if user.Id.Hex() != context.Value("user_id") {
				continue
			}
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

	user.Username = strings.ToLower(user.Username)
	existingUser, err := GetUserByUsername(ctx, user.Username)
	if existingUser != nil {
		return primitive.NilObjectID, ErrDuplicateUsername
	}
	if err != nil && err != ErrUserNotFound {
		return primitive.NilObjectID, err
	}

	now := time.Now()
	user.Id = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Password = HashAndSalt(user.Password)
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
		if err == mongo.ErrNoDocuments {
			return false, ErrUserNotFound
		}
		return false, err
	}

	filter := bson.M{"_id": common.ObjectIDFromHex(id)}
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
		user.Password = HashAndSalt(user.Password)
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
		if err == mongo.ErrNoDocuments {
			return false, ErrUserNotFound
		}
		return false, err
	}

	filter := bson.M{"_id": common.ObjectIDFromHex(id)}

	_, err = db.Collection("users").DeleteOne(*ctx, filter)
	if err != nil {
		return false, err
	}

	return true, nil
}
