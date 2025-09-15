package main

import (
	"Backend/internal/bootstrap"
	"Backend/internal/config"
	"Backend/internal/router"
	"log"
)

func main() {
	// 1. Get env
	config.InitConfig()

	// 2. Init db
	db := bootstrap.NewDB(bootstrap.DBConfig{
		DSN:             config.MYSQL_DSN,
		MaxOpenConn:     config.MYSQL_MAX_OPEN_CONNECTION,
		MaxIdleConn:     config.MYSQL_MAX_IDLE_CONNECTION,
		ConnMaxLifetime: config.MYSQL_CONNECTION_MAX_LIFE_TIME,
		ConnMaxIdleTime: config.MYSQL_CONNECTION_MAX_IDLE_TIME,
	})
	defer func() { _ = db.Close() }()
	rdb := bootstrap.NewRedis(bootstrap.RedisConfig{
		Addr:         config.REDIS_ADDR,
		Password:     config.REDIS_PASSWORD,
		DB:           config.REDIS_DB,
		DialTimeout:  config.REDIS_DIAL_TIMEOUT,
		ReadTimeout:  config.REDIS_READ_TIMEOUT,
		WriteTimeout: config.REDIS_WIRTE_TIMEOUT,
		PingTimeout:  config.REDIS_PING_TIMEOUT,
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
