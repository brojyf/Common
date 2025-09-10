package config

import "fmt"

func RedisKeyThrottle(email, scene string) string {
	return fmt.Sprintf("%s:%s:throttle", email, scene)
}
