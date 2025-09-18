package main

import (
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/router"
	"log"
)

func main() {
	// 1. Init Config
	config.Init()

	// 2. Init DB
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

	// 3. Setup Router: gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(router.Deps{
		DB:  db,
		RDB: rdb,
	})

	// 4. Start Server
	addr := ":8080"
	log.Printf("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server", err)
	}
}
