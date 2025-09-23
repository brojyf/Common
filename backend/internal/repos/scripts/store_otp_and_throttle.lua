-- KEYS[1]=otpKey, KEYS[2]=thKey
-- ARGV[1]=code, ARGV[2]=otpTTL(sec), ARGV[3]=thTTL(sec)

-- 1. 尝试设置 throttle；若已存在则限流
local ok = redis.call("SET", KEYS[2], "1", "NX", "EX", ARGV[3])
if not ok then
    return "THROTTLED"
end

-- 2. 写 OTP
redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])

return "OK"