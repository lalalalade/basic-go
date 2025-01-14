-- 验证码在Redis上的key
local key = KEYS[1]
-- 验证次数
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在但是没有过期时间
    -- 系统错误，没有给过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cnt, 600)
else
    -- 发送太频繁
    return -1
end