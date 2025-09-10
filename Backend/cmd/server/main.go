package main

import (
	"Backend/internal/bootstrap"
	"Backend/internal/config"
	"Backend/internal/router"
	"log"
	"time"
)

func main() {
	// 1. Get env
	config.InitConfig()

	// 2. Init db
	db := bootstrap.NewDB(bootstrap.DBConfig{
		DSN:             config.MYSQL_DSN,
		MaxOpenConn:     30,
		MaxIdleConn:     10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	})
	defer func() { _ = db.Close() }()

	// 3. Set up router
	// Production mode：gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(router.Deps{
		DB: db,
	})

	// 4. Start Server
	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server", err)
	}
}
