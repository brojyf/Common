package scripts

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed store_otp_and_throttle.lua
var storeCodeAndThrottleLua string

type Registry struct {
	StoreCodeAndThrottle *redis.Script
}

func NewRegistry() *Registry {
	return &Registry{
		StoreCodeAndThrottle: redis.NewScript(storeCodeAndThrottleLua),
	}
}
