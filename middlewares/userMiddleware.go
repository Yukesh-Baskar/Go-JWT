package middlewares

import (
	"fmt"
	"go-jwt/models"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

var validate = validator.New()

func RegisterUserMiddleware(c *gin.Context) {
	var user *models.User
	errMsgs := []models.ErrorHandler{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Invalid user name",
			StatusCode: http.StatusBadRequest,
		})
	}

	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    fmt.Sprintf("Invalid user email %v", err.Error()),
			StatusCode: http.StatusBadRequest,
		})
	}

	if user.Gender == "" {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Invalid user gender",
			StatusCode: http.StatusBadRequest,
		})
	}

	if strings.TrimSpace(user.Password) == "" {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Invalid user password",
			StatusCode: http.StatusBadRequest,
		})
	}

	if strings.TrimSpace(user.ConfirmPassword) == "" {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Invalid user confirm password",
			StatusCode: http.StatusBadRequest,
		})
	}

	if strings.TrimSpace(user.Password) != strings.TrimSpace(user.ConfirmPassword) {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Password and Confirm password should be same",
			StatusCode: http.StatusBadRequest,
		})
	}

	if user.Mobile == "" {
		errMsgs = append(errMsgs, models.ErrorHandler{
			Message:    "Invalid user mobile",
			StatusCode: http.StatusBadRequest,
		})
	}

	if len(errMsgs) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, errMsgs)
		return
	}
	c.Set("user", user)
	c.Next()
}

func LoginUserMiddleware(c *gin.Context) {
	var user *models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "Invalid user email or password",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	c.Set("user", user)
	c.Next()
}

func GetUserMiddleWare(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Id not found",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Missing or invalid Authorization header",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	token = token[7:]
	c.Set("bearerToken", token)
	c.Set("id", id)
	c.Next()
}

func UpdateUserMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Missing or invalid Authorization header",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	token = token[7:]

	claims := models.JWTClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY@12345"), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	if !jwtToken.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Invalid Token",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	var user models.User
	err = c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    err.Error(),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	c.Set("claims_user_email", claims.UserEmail)
	c.Set("user", user)
	c.Next()
}

func DeleteUserMiddleware(c *gin.Context) {

}
