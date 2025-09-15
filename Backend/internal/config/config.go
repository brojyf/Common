package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	REDIS_ADDR          string
	REDIS_PASSWORD      string
	REDIS_DB            int
	REDIS_DIAL_TIMEOUT  time.Duration
	REDIS_READ_TIMEOUT  time.Duration
	REDIS_WIRTE_TIMEOUT time.Duration
	REDIS_PING_TIMEOUT  time.Duration

	MYSQL_DSN                      string
	MYSQL_MAX_OPEN_CONNECTION      int
	MYSQL_MAX_IDLE_CONNECTION      int
	MYSQL_CONNECTION_MAX_LIFE_TIME time.Duration
	MYSQL_CONNECTION_MAX_IDLE_TIME time.Duration

	REQUEST_TIMEOUT     time.Duration
	DB_QUERY_TIMEOUT    time.Duration
	REDIS_QUERY_TIMEOUT time.Duration

	OTP_THROTTLE_TTL time.Duration
	OTP_TTL          time.Duration
)

func InitConfig() {
	_ = godotenv.Load("dev.env")

	OTP_THROTTLE_TTL = mustGetEnvAsDuration("OTP_THROTTLE_TTL")
	OTP_TTL = mustGetEnvAsDuration("OTP_TTL")

	REQUEST_TIMEOUT = mustGetEnvAsDuration("REQUEST_TIMEOUT")
	DB_QUERY_TIMEOUT = mustGetEnvAsDuration("DB_QUERY_TIMEOUT")
	REDIS_QUERY_TIMEOUT = mustGetEnvAsDuration("REDIS_QUERY_TIMEOUT")

	REDIS_ADDR = mustGetEnv("REDIS_ADDR")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	REDIS_DB = mustGetEnvAsInt("REDIS_DB")
	REDIS_DIAL_TIMEOUT = mustGetEnvAsDuration("REDIS_DIAL_TIMEOUT")
	REDIS_READ_TIMEOUT = mustGetEnvAsDuration("REDIS_READ_TIMEOUT")
	REDIS_WIRTE_TIMEOUT = mustGetEnvAsDuration("REDIS_WIRTE_TIMEOUT")
	REDIS_PING_TIMEOUT = mustGetEnvAsDuration("REDIS_PING_TIMEOUT")

	MYSQL_DSN = mustGetEnv("MYSQL_DSN")
	MYSQL_MAX_OPEN_CONNECTION = mustGetEnvAsInt("MYSQL_MAX_OPEN_CONNECTION")
	MYSQL_MAX_IDLE_CONNECTION = mustGetEnvAsInt("MYSQL_MAX_IDLE_CONNECTION")
	MYSQL_CONNECTION_MAX_LIFE_TIME = mustGetEnvAsDuration("MYSQL_CONNECTION_MAX_LIFE_TIME")
	MYSQL_CONNECTION_MAX_IDLE_TIME = mustGetEnvAsDuration("MYSQL_CONNECTION_MAX_IDLE_TIME")
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
