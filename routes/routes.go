package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/auth/register", controllers.Register)

	return router
}
