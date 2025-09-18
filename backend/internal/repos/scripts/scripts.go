package scripts

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

type Registry struct {
	StoreOTPAndThrottle         *redis.Script
	ThrottleMatchAndConsumeCode *redis.Script
}

func NewRegistry() *Registry {
	return &Registry{
		StoreOTPAndThrottle:         redis.NewScript(storeOTPAndThrottleLua),
		ThrottleMatchAndConsumeCode: redis.NewScript(throttleMatchAndConsumeCodeLua),
	}
}

//go:embed throttle_match_consume_code.lua
var throttleMatchAndConsumeCodeLua string

//go:embed store_otp_and_throttle.lua
var storeOTPAndThrottleLua string
