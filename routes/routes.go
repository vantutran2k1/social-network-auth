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

	router.PUT("/api/auth/password", middlewares.AuthMiddleware(), controllers.UpdatePassword)

	router.POST("/api/auth/password-reset-token", controllers.CreateResetPasswordToken)
	router.POST("/api/auth/password-reset-email", controllers.SendResetPasswordEmail)
	router.PUT("/api/auth/password-reset", controllers.ResetPassword)

	router.PATCH("/api/auth/level", controllers.UpdateUserLevel)

	router.GET("/api/profiles", controllers.GetProfile)
	router.GET("/api/profiles/me", middlewares.AuthMiddleware(), controllers.GetCurrentProfile)
	router.POST("/api/profiles", middlewares.AuthMiddleware(), controllers.CreateProfile)
	router.PUT("/api/profiles", middlewares.AuthMiddleware(), controllers.UpdateCurrentProfile)

	return router
}
