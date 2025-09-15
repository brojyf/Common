package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Redis struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PingTimeout  time.Duration
}

type MySQL struct {
	DSN                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

type Timeouts struct {
	Request     time.Duration
	RequestCode time.Duration
}

type RedisTTL struct {
	OTP         time.Duration
	OTPThrottle time.Duration
}

type Config struct {
	Redis    Redis
	MySQL    MySQL
	Timeouts Timeouts
	RedisTTL RedisTTL
}

var C Config

func Init() {
	_ = godotenv.Load("dev.env")

	C = Config{
		Redis: Redis{
			Addr:         mustGet("REDIS_ADDR"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           mustGetInt("REDIS_DB"),
			DialTimeout:  mustGetDur("REDIS_DIAL_TIMEOUT"),
			ReadTimeout:  mustGetDur("REDIS_READ_TIMEOUT"),
			WriteTimeout: mustGetDur("REDIS_WRITE_TIMEOUT"),
			PingTimeout:  mustGetDur("REDIS_PING_TIMEOUT"),
		},
		MySQL: MySQL{
			DSN:                   mustGet("MYSQL_DSN"),
			MaxOpenConnections:    mustGetInt("MYSQL_MAX_OPEN_CONNECTION"),
			MaxIdleConnections:    mustGetInt("MYSQL_MAX_IDLE_CONNECTION"),
			ConnectionMaxLifetime: mustGetDur("MYSQL_CONNECTION_MAX_LIFE_TIME"),
			ConnectionMaxIdleTime: mustGetDur("MYSQL_CONNECTION_MAX_IDLE_TIME"),
		},
		Timeouts: Timeouts{
			Request:     mustGetDur("REQUEST_TIMEOUT"),
			RequestCode: mustGetDur("REQUEST_CODE"),
		},
		RedisTTL: RedisTTL{
			OTP:         mustGetDur("OTP_TTL"),
			OTPThrottle: mustGetDur("OTP_THROTTLE_TTL"),
		},
	}

}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Required environment variable %s is missing", key)
	}
	return val
}

func mustGetDur(key string) time.Duration {
	val := mustGet(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("❌ Invalid duration for %s: %v", key, err)
	}
	return d
}

func mustGetInt(key string) int {
	val := mustGet(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("❌ Invalid int for %s: %v", key, err)
	}
	return num
}
