package main

import (
	"Backend/internal/bootstrap"
	"Backend/internal/config"
	"Backend/internal/router"
	"log"
)

func main() {
	// 1. Get env
	config.Init()

	// 2. Init db
	db := bootstrap.NewDB(bootstrap.DBConfig{
		DSN:             config.C.MySQL.DSN,
		MaxOpenConn:     config.C.MySQL.MaxOpenConnections,
		MaxIdleConn:     config.C.MySQL.MaxIdleConnections,
		ConnMaxLifetime: config.C.MySQL.ConnectionMaxLifetime,
		ConnMaxIdleTime: config.C.MySQL.ConnectionMaxIdleTime,
	})
	defer func() { _ = db.Close() }()
	rdb := bootstrap.NewRedis(bootstrap.RedisConfig{
		Addr:         config.C.Redis.Addr,
		Password:     config.C.Redis.Password,
		DB:           config.C.Redis.DB,
		DialTimeout:  config.C.Redis.DialTimeout,
		ReadTimeout:  config.C.Redis.ReadTimeout,
		WriteTimeout: config.C.Redis.WriteTimeout,
		PingTimeout:  config.C.Redis.PingTimeout,
	})
	defer func() { _ = rdb.Close() }()

	// 3. Set up router
	// Production mode：gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(router.Deps{
		DB:  db,
		RDB: rdb,
	})

	// 4. Start Server
	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server", err)
	}
}
