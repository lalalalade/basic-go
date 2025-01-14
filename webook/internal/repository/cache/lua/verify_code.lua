local key = KEYS[1]

local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get",cntKey))
if cnt <= 0 then
    --    用户一直输错 或者 已经用过了
    return -1
elseif expectedCode == code then
    redis.call("set", cntKey, -1)
    return 0
else
--    输错了
    redis.call("decr", cntKey, -1)
    return -2
end