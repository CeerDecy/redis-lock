package lock

const (
	// 尝试上锁
	lockScript = `
local val = redis.call('get', KEYS[1])
if val == false then
    -- key 不存在
    return redis.call('set', KEYS[1], ARGV[1], 'EX', ARGV[2])
elseif val == ARGV[1] then
    redis.call('expire', KEYS[1], ARGV[2])
    return  "OK"
else
    return ""
end`
	// 刷新锁的存货时间
	refreshScript = `
if redis.call("get", KEYS[1]) == ARGV[1]
then
    return redis.call("expire", KEYS[1], ARGV[2])
else
    return 0
end`

	// 尝试解锁
	unlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1]
then
    return redis.call("del", KEYS[1])
else
    return 0
end`
)
