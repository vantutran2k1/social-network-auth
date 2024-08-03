package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/controllers"
	"github.com/vantutran2k1/social-network-auth/middlewares"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)
	router.POST("/api/auth/logout", middlewares.AuthMiddleware(), controllers.Logout)

	router.GET("/api/levels", controllers.GetLevels)
	router.POST("/api/levels/assign", controllers.AssignLevelToUser)

	return router
}
