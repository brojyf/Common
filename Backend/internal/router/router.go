package router

import (
	"Backend/internal/config"
	"Backend/internal/handlers"
	"Backend/internal/middlewares"
	"Backend/internal/repos"
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
		}
	}

	// 5. Return router
	return r
}

func checkHealth(c *gin.Context) {
	c.String(200, "pong")
}
