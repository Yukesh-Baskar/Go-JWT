package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id              primitive.ObjectID `bson:"_id"`
	FirstName       string             `json:"first_name" validate:"required,min=2,max=200"`
	LastName        string             `json:"last_name" validate:"required,min=2,max=200"`
	Email           string             `json:"email" validate:"required"`
	Age             int                `json:"age" validate:"required"`
	Gender          string             `json:"gender" validate:"required,eq=Male|eq=Female"`
	Mobile          string             `json:"mobile" validate:"required"`
	Password        string             `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string             `json:"confirm_password" validate:"required,min=6,max=100"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
}

type JWTClaims struct {
	Username  string
	UserEmail string
	Mobile    string
	jwt.RegisteredClaims
}
