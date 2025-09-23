package scripts

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

type Registry struct {
	StoreOTPAndThrottle         *redis.Script
	ThrottleMatchAndConsumeCode *redis.Script
	FindAdnMarkOTTJTI           *redis.Script
}

func NewRegistry() *Registry {
	return &Registry{
		StoreOTPAndThrottle:         redis.NewScript(storeOTPAndThrottleLua),
		ThrottleMatchAndConsumeCode: redis.NewScript(throttleMatchAndConsumeCodeLua),
		FindAdnMarkOTTJTI:           redis.NewScript(findAndMarkOTTJTI),
	}
}

//go:embed find_jti_usage_and_mark.lua
var findAndMarkOTTJTI string

//go:embed throttle_match_consume_code.lua
var throttleMatchAndConsumeCodeLua string

//go:embed store_otp_and_throttle.lua
var storeOTPAndThrottleLua string
