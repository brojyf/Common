-- KEYS[1] = otp key
-- KEYS[2] = throttle key
-- ARGV[1] = code
-- ARGV[2] = otp ttl (秒)
-- ARGV[3] = throttle ttl (秒)

redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
redis.call("SET", KEYS[2], "1", "EX", ARGV[3])
return "OK"