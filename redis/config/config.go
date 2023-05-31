package config

import (
	"github.com/magiconair/properties"
	"log"
	"os"
	"sync"
)

// instance 用于获取Config中的参数
var instance *RedisConfig

// RedisConfigInstance 单例模式获取
func RedisConfigInstance() *RedisConfig {
	once := sync.Once{}
	once.Do(func() {
		instance = &RedisConfig{
			properties: initProperties(),
		}
	})
	return instance
}

// 不使用单例模式，也可以通过init初始化
//func init() {
//	instance = &RedisConfig{
//		properties: initProperties(),
//	}
//}

func initProperties() *properties.Properties {
	basePath, _ := os.Getwd()
	p, err := properties.LoadFile(basePath+"/config/redis.properties", properties.UTF8)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

type RedisConfig struct {
	properties *properties.Properties
}

func (config *RedisConfig) GetHost() string {
	return config.properties.MustGet("host")
}

func (config *RedisConfig) GetPort() string {
	return config.properties.MustGet("port")
}
