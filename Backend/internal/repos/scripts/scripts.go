package scripts

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed store_otp_and_throttle.lua
var storeOTPAndThrottleLua string

type Registry struct {
	StoreOTPAndThrottle *redis.Script
}

func NewRegistry() *Registry {
	return &Registry{
		StoreOTPAndThrottle: redis.NewScript(storeOTPAndThrottleLua),
	}
}
