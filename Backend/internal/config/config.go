package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	REDIS_ADDR     string
	REDIS_PASSWORD string
	MYSQL_DSN      string

	REQUEST_TIMEOUT     time.Duration
	DB_QUERY_TIMEOUT    time.Duration
	REDIS_QUERY_TIMEOUT time.Duration
)

func InitConfig() {
	_ = godotenv.Load("dev.env")

	REDIS_ADDR = mustGetEnv("REDIS_ADDR")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	MYSQL_DSN = mustGetEnv("MYSQL_DSN")

	REQUEST_TIMEOUT = mustGetEnvAsDuration("REQUEST_TIMEOUT")
	DB_QUERY_TIMEOUT = mustGetEnvAsDuration("DB_QUERY_TIMEOUT")
	REDIS_QUERY_TIMEOUT = mustGetEnvAsDuration("REDIS_QUERY_TIMEOUT")
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Required environment variable %s is missing", key)
	}
	return val
}

func mustGetEnvAsDuration(key string) time.Duration {
	val := mustGetEnv(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("❌ Invalid duration for %s: %v", key, err)
	}
	return d
}

func mustGetEnvAsInt(key string) int {
	val := mustGetEnv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("❌ Invalid int for %s: %v", key, err)
	}
	return num
}
