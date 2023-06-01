package lock

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"redis-lock/redis/client"
	"time"
)

// RedisLock 分布式锁对象
type RedisLock struct {
	uuid   string             // uuid记录锁的所有者
	client *client.Client     // Redis的连接对象
	cancel context.CancelFunc // 用于取消子Goroutine的更新操作
	state  bool               // 当前锁的状态
	key    string             // 锁的名称，即Redis中的key
}

// MakeRedisLock 获取分布式锁对象
func MakeRedisLock(ctx context.Context, key string) *RedisLock {
	return &RedisLock{
		uuid:   uuid.New().String(),
		client: client.GetClient(ctx),
		state:  false,
		key:    key,
	}
}

// 启动子Goroutine用于更新Redis中key的存活时间
func (lock *RedisLock) refresh(ctx context.Context) {
	// 循环监听ctx是否需要关闭当前Goroutine
	for {
		select {
		// Channel中有数据则关闭当前Goroutine
		case <-ctx.Done():
			return
		default:
			// 否则执行更新操作
			res, err := lock.client.RunLua(refreshScript, []string{lock.key}, []any{lock.uuid, 5})
			// 若更新失败则直接放弃该锁，关闭当前Goroutine
			if err != nil || res != int64(1) {
				return
			}
			// 每三秒更新一次
			time.Sleep(3 * time.Second)
		}
	}
}

// TryLock 尝试获取锁
func (lock *RedisLock) TryLock() (bool, error) {
	// 尝试获取锁，这里的查询锁的拥有者和执行加锁是原子化的操作，因此需要写在lua脚本里执行
	res, err := lock.client.RunLua(lockScript, []string{lock.key}, lock.uuid, 5)
	// 获取失败则向上层应该抛出异常
	if err != nil {
		return false, err
	}
	// 若res返回OK说明上锁成功
	if res == "OK" {
		// 设置带Cancel函数的context
		ctx, cancelFunc := context.WithCancel(context.Background())
		lock.cancel = cancelFunc
		// 上锁状态设置为true
		lock.state = true
		// 开启协程不断更新key的存活时间
		go lock.refresh(ctx)
		return true, nil
	}
	return false, nil
}

// UnLock 尝试解开分布式锁
func (lock *RedisLock) UnLock() (bool, error) {
	// 若锁的状态为false,说明该锁还没有被上锁无法进行解锁操作
	if !lock.state {
		return false, errors.New("this lock [" + lock.key + "] is not locked")
	}
	// 执行Lua脚本
	res, err := lock.client.RunLua(unlockScript, []string{lock.key}, lock.uuid)
	// 不管成功与否都要讲state改为false
	lock.state = false
	if err != nil {
		return false, err
	}
	// 若成功解锁，那么关闭子Goroutine的更新操作
	if res == int64(1) {
		lock.cancel()
		return true, nil
	}
	return false, nil
}
