package main

import (
	"Backend/internal/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Get env
	if err := godotenv.Load("dev.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Init db

	// 3. Set up router
	// Production mode：gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter()

	// 4. Start Server
	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server", err)
	}
}
