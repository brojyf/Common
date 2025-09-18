package config

import (
	"fmt"
)

// RedisKeyOTP opt:<email>:<scene>:<codeID>
func RedisKeyOTP(email, scene, codeID string) string {
	return fmt.Sprintf("otp:%s:%s:%s", email, scene, codeID)
}

// RedisKeyThrottle opt:throttle:<email>:<scene>
func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("otp:throttle:%s:%s", email, scene)
}
