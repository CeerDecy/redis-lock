package lock

import (
	"context"
	"github.com/google/uuid"
	"log"
	"redis-lock/redis/client"
	"time"
)

const (
	lockScript = `
local val = redis.call('get', KEYS[1])
-- 在加锁的重试的时候，要判断自己上一次是不是加锁成功了
if val == false then
    -- key 不存在
    return redis.call('set', KEYS[1], ARGV[1], 'EX', ARGV[2])
elseif val == ARGV[1] then
    -- 刷新过期时间
    redis.call('expire', KEYS[1], ARGV[2])
    return  "OK"
else
    -- 此时别人持有锁
    return ""
end`
)

type RedisLock struct {
	uuid   string
	client *client.Client
	cancel context.CancelFunc
	state  bool
}

func MakeRedisLock() *RedisLock {
	ctx := context.Background()
	return &RedisLock{
		uuid:   uuid.New().String(),
		client: client.GetClient(ctx),
		state:  false,
	}
}

func (lock *RedisLock) refresh(ctx context.Context, key string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(3 * time.Second)
			lua, err := lock.client.RunLua(lockScript, []string{key}, []any{lock.uuid, 5})
			if err != nil {
				log.Fatal(err)
			}
			log.Println("refresh key ", key, " ", lua)
		}
	}
}

func (lock *RedisLock) TryLock(key string) (bool, error) {
	res, err := lock.client.RunLua(lockScript, []string{key}, []any{lock.uuid, 5})
	if err != nil {
		return false, err
	}
	if res == "OK" {
		ctx, cancelFunc := context.WithCancel(context.Background())
		lock.cancel = cancelFunc
		lock.state = true
		go lock.refresh(ctx, key)
		return true, nil
	}
	return false, nil
}

func (lock *RedisLock) UnLock() {
	if lock.state {
		lock.cancel()
	}
}
