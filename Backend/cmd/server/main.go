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

	// Set up router
	// Production mode：gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter()

	// Start Server
	addr := ":8080"
	log.Printf("Listening on %s...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
