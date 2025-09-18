package config

import (
	"fmt"
)

// RedisKeyVerifyThrottle verify:throttle:<email>:<scene>
func RedisKeyVerifyThrottle(email, scene string) string {
	return fmt.Sprintf("verify:throttle:%s:%s", email, scene)
}

// RedisKeyOTP otp:<email>:<scene>:<codeID>
func RedisKeyOTP(email, scene, codeID string) string {
	return fmt.Sprintf("otp:%s:%s:%s", email, scene, codeID)
}

// RedisKeyThrottle otp:throttle:<email>:<scene>
func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
