package main

import (
	"context"
	"fmt"
	"redis-lock/lock"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	redisLock := lock.MakeRedisLock(context.Background(), "service")
	tryLock, err := redisLock.TryLock()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tryLock)
	if !tryLock {
		return
	}
	time.Sleep(20 * time.Second)
	_, err = redisLock.UnLock()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}
