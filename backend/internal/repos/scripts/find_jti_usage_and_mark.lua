-- Return 0 if okay, 1 otherwise

-- 1.1 Find jti
local val = redis.call("GET", KEYS[1])
if val == "0" then
    -- key 存在并且是 0 -> 更新为 1
    redis.call("SET", KEYS[1], 1, "EX", tonumber(ARGV[1]))
    return 0  -- success
elseif val == "1" then
    return 1
else
    return 1
end