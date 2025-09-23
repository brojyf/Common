-- 固定窗口：检查 -> 写入计数/过期 -> 校验并消费 OTP
-- KEYS[1]=rl_counter_key, KEYS[2]=otp_key
-- ARGV[1]=limit, ARGV[2]=window_sec, ARGV[3]=expected_code

local limit      = tonumber(ARGV[1])
local window_sec = tonumber(ARGV[2])
local expected   = ARGV[3]

-- 1.1 先检查是否超限
local cur = tonumber(redis.call("GET", KEYS[1]) or "0")
if cur >= limit then
    local ttl = redis.call("TTL", KEYS[1]) or -1
    return "THROTTLED"
end

-- 1.2 写 throttle：INCR 并设置窗口过期（第一次创建时）
cur = redis.call("INCR", KEYS[1])
if cur == 1 then
    redis.call("EXPIRE", KEYS[1], window_sec)
end

-- 1.3 如果写入后超过 limit，立即返回限流
if cur > limit then
    local ttl = redis.call("TTL", KEYS[1]) or -1
    return "THROTTLED"
end

-- 2. 找 code
local v = redis.call("GET", KEYS[2])
if not v then
    return "EXPIRED"
end
if v ~= expected then
    return "INVALID"
end

-- 3. 匹配则消费（删除）
redis.call("DEL", KEYS[2])

-- 4. 设置JTI有效
redis.call("SETEX", KEYS[3], tonumber(ARGV[4]), 0)
return "OK"
