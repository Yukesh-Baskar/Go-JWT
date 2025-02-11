package routes

import (
	"go-jwt/controllers"
	"go-jwt/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoute *gin.Engine) {
	// sign-up
	incomingRoute.POST("/signup", middlewares.RegisterUserMiddleware, controllers.RegisterUser)

	// login
	incomingRoute.POST("/login", middlewares.LoginUserMiddleware, controllers.LoginUser)

	// get-user
	incomingRoute.GET("/user/:id", middlewares.GetUserMiddleWare, controllers.GetUserController)

	// update-user
	incomingRoute.PATCH("/update-user/:id", middlewares.UpdateUserMiddleware, controllers.UpdateUserController)

	// delete-user
	incomingRoute.DELETE("/delete-user/:id", middlewares.DeleteUserMiddleware, controllers.DeleteUserController)
}
