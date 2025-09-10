package router

import (
	"Backend/internal/handlers"
	"Backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ping", checkHealth)

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/refresh", handlers.HandleRefresh)
			authGroup.POST("/login", handlers.HandleLogin)
			authGroup.POST("/request-code", handlers.HandleRequestCode)
			authGroup.POST("/verify-code", handlers.HandleVerifyCode)
			authGroup.POST("/create-account", middleware.OPTMiddleware(), handlers.HandleCreateAccount)
			authGroup.POST("/forget-password", middleware.OPTMiddleware(), handlers.HandleForgetPassword)
			authGroup.PATCH("/reset-password", middleware.ATKMiddleware(), handlers.HandleResetPassword)
			authGroup.POST("/logout-all", middleware.ATKMiddleware(), handlers.HandleLogoutAll)
			authGroup.PATCH("/me/set-username", middleware.ATKMiddleware(), handlers.HandleSetUsername)
		}
	}

	return r
}

func checkHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
