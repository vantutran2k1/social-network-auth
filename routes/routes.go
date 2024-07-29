package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)
	router.POST("/api/auth/logout", controllers.Logout)
	router.POST("/api/auth/validate", controllers.Validate)

	return router
}
