package router

import (
	"Backend/internal/config"
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

	// context: RID -> Timeout
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Timeout(config.C.Timeouts.Request))
	r.Use(middleware.AccessLog())

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
			authGroup.POST("/create-account", middleware.OPT(), authH.HandleCreateAccount)
			authGroup.POST("/forget-password", middleware.OPT(), authH.HandleForgetPassword)
			authGroup.PATCH("/reset-password", middleware.ATK(), authH.HandleResetPassword)
			authGroup.POST("/logout-all", middleware.ATK(), authH.HandleLogoutAll)
			authGroup.PATCH("/me/set-username", middleware.ATK(), authH.HandleSetUsername)
			authGroup.POST("/logout", middleware.ATK(), authH.HandleLogout)
		}
	}

	return r
}

func checkHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
