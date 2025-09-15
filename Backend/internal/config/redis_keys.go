package config

import (
	"fmt"
)

func RedisKeyOTP(email, scene, jti string) string {
	return fmt.Sprintf("otp:%s:%s:%s", email, scene, jti)
}

func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
