package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	REDIS_ADDR     string
	REDIS_PASSWORD string

	MYSQL_DSN string
)

func InitConfig() {
	_ = godotenv.Load("dev.env")

	REDIS_ADDR = mustGetEnv("REDIS_ADDR")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	MYSQL_DSN = mustGetEnv("MYSQL_DSN")
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Required environment variable %s is missing", key)
	}
	return val
}

func mustGetEnvAsInt(key string) int {
	val := mustGetEnv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("❌ Invalid int for %s: %v", key, err)
	}
	return num
}
