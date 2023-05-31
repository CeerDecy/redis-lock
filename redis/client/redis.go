package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
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
func (client *Client) SetNX(key string, value any) bool {
	result, err := client.rdb.SetNX(client.ctx, key, value, 5*time.Second).Result()
	if err != nil {
		log.Println("set key error :", err)
	}
	log.Println("set key ", result)
	return result
}

// SetWithTimeNX 设置Key和Value
func (client *Client) SetWithTimeNX(key string, value any, t time.Duration) bool {
	result, err := client.rdb.SetNX(client.ctx, key, value, t).Result()
	if err != nil {
		log.Println("set key error :", err)
	}
	log.Println("set key ", result)
	return result
}

// Get 获取当前Key的值
func (client *Client) Get(key string) string {
	result, err := client.rdb.Get(client.ctx, key).Result()
	if err != nil {
		log.Println("set key error :", err)
	}
	log.Println("set key ", result)
	return result
}

// DeleteKey 删除Key
func (client *Client) DeleteKey(key string) int {
	result, err := client.rdb.Del(client.ctx, key).Result()
	if err != nil {
		log.Println("delete key error ", err)
	}
	log.Println("del key ", result)
	return int(result)
}

func (client *Client) RunLua(script string, keys []string, args ...any) {
	client.rdb.Eval(client.ctx, script, keys, args...)
}
