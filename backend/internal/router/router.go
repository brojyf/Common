package router

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middlewares"
	"backend/internal/repos"
	"backend/internal/services"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	DB  *sql.DB
	RDB *redis.Client
}

func SetupRouter(d Deps) *gin.Engine {
	// 1. Set Up Engine
	r := gin.Default()

	// 2. User Middlewares
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestID())
	r.Use(middlewares.Timeout(config.C.Timeouts.Request))
	r.Use(middlewares.AccessLog())

	// 3. Dependencies Injection
	authRepo := repos.NewAuthRepo(d.DB, d.RDB)
	authSvc := services.NewAuthService(authRepo)
	authH := handlers.NewAuthHandler(authSvc)

	// 4. Register Router
	apiGroup := r.Group("/api")
	{
		// a. Health Check
		apiGroup.GET("/ping", checkHealth)

		// b. Auth
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/request-code", authH.HandleRequestCode)
			authGroup.POST("/verify-code", authH.HandleVerifyCode)
			authGroup.POST("/create-account", middlewares.OneTimeToken(), authH.HandleCreateAccount)
			authGroup.POST("/login", authH.HandleLogin)
		}
	}

	// 5. Return router
	return r
}

func checkHealth(c *gin.Context) {
	c.String(200, "pong")
}
