package router

import (
	"Backend/internal/handlers"
	"Backend/internal/middleware"
	"Backend/internal/repo"
	"Backend/internal/services"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	DB  *sql.DB
	RDB *redis.Client
}

func SetupRouter(d Deps) *gin.Engine {
	r := gin.Default()

	authRepo := repo.NewAuthRepo(d.DB, d.RDB)
	authSvc := services.NewAuthService(authRepo)
	authH := handlers.NewAuthHandler(authSvc)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ping", checkHealth)

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/refresh", authH.HandleRefresh)
			authGroup.POST("/login", authH.HandleLogin)
			authGroup.POST("/request-code", authH.HandleRequestCode)
			authGroup.POST("/verify-code", authH.HandleVerifyCode)
			authGroup.POST("/create-account", middleware.OPTMiddleware(), authH.HandleCreateAccount)
			authGroup.POST("/forget-password", middleware.OPTMiddleware(), authH.HandleForgetPassword)
			authGroup.PATCH("/reset-password", middleware.ATKMiddleware(), authH.HandleResetPassword)
			authGroup.POST("/logout-all", middleware.ATKMiddleware(), authH.HandleLogoutAll)
			authGroup.PATCH("/me/set-username", middleware.ATKMiddleware(), authH.HandleSetUsername)
			authGroup.POST("/logout", middleware.ATKMiddleware(), authH.HandleLogout)
		}
	}

	return r
}

func checkHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
