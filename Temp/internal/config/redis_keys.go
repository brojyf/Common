package config

import (
	"fmt"
)

func RedisKeyOTP(email, scene, codeID string) string {
	return fmt.Sprintf("otp:%s:%s:%s", email, scene, codeID)
}

func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
