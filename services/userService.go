package services

import (
	"context"
	"fmt"
	"go-jwt/database"
	"go-jwt/models"
	"go-jwt/utils"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection = database.OpenColletion(database.Client, "user")

func Signup(user *models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}, {Key: "mobile", Value: user.Mobile}}
	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if count > 0 {
		return nil, &models.ErrorHandler{
			Message:    "User with this email or mobile already exists",
			StatusCode: http.StatusBadRequest,
		}
	}
	user.Id = primitive.NewObjectID()
	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	hPass, err := utils.GetHashedPassword(user.Password)
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	user.Password = string(hPass)
	user.ConfirmPassword = string(hPass)
	// hash the password -> skipping as of now

	res, err := userCollection.InsertOne(ctx, &user)
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    fmt.Sprintf("error occured while inserting document: %v", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return res, nil
}

func LoginUser(user *models.User) (token string, err error) {
	var dbUser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}}
	if err := userCollection.FindOne(ctx, filter).Decode(&dbUser); err != nil {
		return token, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}
	fmt.Println(dbUser.Password, user.Password)
	err = utils.CheckPassword(dbUser.Password, user.Password)
	if err != nil {
		fmt.Println(err)
		return token, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}
	// if user.Password != dbUser.Password {
	// 	return token, &models.ErrorHandler{
	// 		Message:    "Incorrect Email or Password",
	// 		StatusCode: http.StatusUnauthorized,
	// 	}
	// }

	claims := models.JWTClaims{
		Username:  dbUser.FirstName,
		UserEmail: dbUser.Email,
		Mobile:    dbUser.Mobile,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(5 * time.Minute),
			},
		},
	}

	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwttoken.SignedString([]byte("SECRET_KEY@12345"))
}

func GetUserService(id, token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	claims := &models.JWTClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY@12345"), nil
	})
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}
	if !jwtToken.Valid {
		return nil, &models.ErrorHandler{
			Message:    "Not a valid token",
			StatusCode: http.StatusUnauthorized,
		}
	}
	var dbUser models.User
	objectIdFromHex, _ := primitive.ObjectIDFromHex(strings.TrimSpace(id))
	filter := bson.D{{Key: "_id", Value: objectIdFromHex}}
	err = userCollection.FindOne(ctx, filter).Decode(&dbUser)
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	return &dbUser, nil
}

func UpdateUser(claims_user_email string, user models.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: claims_user_email}})
	if err != nil {
		return nil, &models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	if count == 0 {
		return nil, &models.ErrorHandler{
			Message:    "No documents with this email",
			StatusCode: http.StatusBadRequest,
		}
	}

	filter := bson.D{{Key: "email", Value: claims_user_email}}
	updateObj := primitive.D{}

	if user.Age > 0 {
		updateObj = append(updateObj, bson.E{Key: "age", Value: user.Age})
	}
	if user.Email != "" {
		updateObj = append(updateObj, bson.E{Key: "email", Value: user.Email})
	}
	if user.FirstName != "" {
		updateObj = append(updateObj, bson.E{Key: "firstname", Value: user.FirstName})
	}
	if user.LastName != "" {
		updateObj = append(updateObj, bson.E{Key: "lastname", Value: user.LastName})
	}
	if user.Mobile != "" {
		updateObj = append(updateObj, bson.E{Key: "mobile", Value: user.Mobile})
	}
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: user.Updated_At})
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	return userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opts)
}
