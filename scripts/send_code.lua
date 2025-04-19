local key = KEYS[1] -- 存储验证码的key code:biz:phone
local cnt = key..":cnt" -- 存储剩余验证次数的key code:biz:phone:cnt
local val = ARGV[1] -- 发送的验证码
local ttl = tonumber(redis.call("ttl", key)) -- 验证码剩余有效期

if ttl == -1 then -- -1 表示没有设置过期时间
    return -1 -- key 存在 但没有设置过期时间 说明系统异常 用 -1 表示
elseif ttl == -2 or ttl < 840 then -- -2 表示 key 不存在 or 剩余过期时间小于14分钟 用于保证一分钟内只能发送一次验证码
    redis.call("set", key, val)
    redis.call("expire", key, 900)
    redis.call("set", cnt, 3)
    redis.call("expire", cnt, 900)
    return 0 -- 正常发送验证码
else
    return -2 -- 发送频率过高啦
end