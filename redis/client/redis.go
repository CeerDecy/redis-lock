package client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"redis-lock/redis/config"
	"time"
)

var redisClient *Client

// Client Redis客户端
type Client struct {
	rdb *redis.Client
	ctx context.Context
}

// 初始化Redis连接
func init() {
	redisClient = &Client{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisConfigInstance().GetAddress(),
		}),
	}
}

// GetClient 获取客户端
func GetClient(ctx context.Context) *Client {
	redisClient.ctx = ctx
	return redisClient
}

// SetNX 设置Key和Value，默认5秒后删除
func (client *Client) SetNX(key string, value any) (bool, error) {
	return client.rdb.SetNX(client.ctx, key, value, 5*time.Second).Result()
}

// SetWithTimeNX 设置Key和Value
func (client *Client) SetWithTimeNX(key string, value any, t time.Duration) (bool, error) {
	return client.rdb.SetNX(client.ctx, key, value, t).Result()
}

// Get 获取当前Key的值
func (client *Client) Get(key string) (string, error) {
	return client.rdb.Get(client.ctx, key).Result()
}

// DeleteKey 删除Key
func (client *Client) DeleteKey(key string) (int64, error) {
	return client.rdb.Del(client.ctx, key).Result()
}

// RunLua 执行Lua脚本
func (client *Client) RunLua(script string, keys []string, args ...any) (any, error) {
	return client.rdb.Eval(client.ctx, script, keys, args...).Result()
}
