package controllers

import (
	"go-jwt/models"
	"go-jwt/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	cUser, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "User not found in context",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}
	user, ok := cUser.(*models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "User is not of the expected type",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}
	res, err := services.Signup(user)
	if err != nil {
		errs, _ := err.(*models.ErrorHandler)
		c.AbortWithStatusJSON(errs.StatusCode, models.ErrorHandler{
			Message:    errs.Message,
			StatusCode: errs.StatusCode,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func LoginUser(c *gin.Context) {
	cUser, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "User not found in context",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}
	user, ok := cUser.(*models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "User is not of the expected type",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}
	token, err := services.LoginUser(user)
	if err != nil {
		c.AbortWithStatusJSON(err.(*models.ErrorHandler).StatusCode, models.ErrorHandler{
			Message:    err.(*models.ErrorHandler).Message,
			StatusCode: err.(*models.ErrorHandler).StatusCode,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"token":   token,
	})
}

func GetUserController(c *gin.Context) {
	token, ok := c.Get("bearerToken")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "Error occured while fetching token data",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}
	id, ok := c.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorHandler{
			Message:    "Error occured while fetching id data",
			StatusCode: http.StatusInsufficientStorage,
		})
		return
	}

	user, err := services.GetUserService(id.(string), token.(string))
	if err != nil {
		c.AbortWithStatusJSON(err.(*models.ErrorHandler).StatusCode, models.ErrorHandler{
			Message:    err.(*models.ErrorHandler).Message,
			StatusCode: err.(*models.ErrorHandler).StatusCode,
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUserController(c *gin.Context) {
	claims_user_email, ok := c.Get("claims_user_email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Error occured while getting claims data",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorHandler{
			Message:    "Error occured while getting user data",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	res, err := services.UpdateUser(claims_user_email.(string), user.(models.User))
	if err != nil {
		c.AbortWithStatusJSON(err.(*models.ErrorHandler).StatusCode, models.ErrorHandler{
			Message:    err.(*models.ErrorHandler).Message,
			StatusCode: err.(*models.ErrorHandler).StatusCode,
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
