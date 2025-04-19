local key = KEYS[1] -- 存储验证码的key code:biz:phone
local cnt = key..":cnt" -- 存储剩余验证次数的key code:biz:phone:cnt
local val = ARGV[1] -- 用户输入的验证码

if tonumber(redis.call("get", cnt)) <= 0 then
    return -1 -- 验证次数已用完 或者 验证过了
elseif redis.call("get", key) == val then
    redis.call("set", cnt, 0) -- 验证成功后将剩余验证次数置为 0
    return 0 -- 验证码正确
else
    redis.call("decr", cnt)
    return -2 -- 验证码错误
end