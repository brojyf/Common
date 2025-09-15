package config

import (
	"fmt"
	"time"
)

func RedisKeyOTP(email, scene string) string {
	ts := time.Now().Unix()
	return fmt.Sprintf("otp:%s:%s:%d", email, scene, ts)
}

func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
