package app

import (
	"fmt"
	"go-jwt/routes"
	"go-jwt/utils"

	"github.com/gin-gonic/gin"
)

func StartApplication() error {
	port, _, err := utils.GetConfig()
	if err != nil {
		return err
	}
	// gin.ForceConsoleColor()
	router := gin.Default()
	routes.UserRoutes(router)

	fmt.Printf("Server Listening on port: %v", port)
	err = router.Run(port)
	if err != nil {
		return err
	}

	return nil
}
