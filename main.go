package main

import (
	"fmt"
	"redis-lock/lock"
	"time"
)

// git push -u origin main
func main() {
	redisLock := lock.MakeRedisLock()
	tryLock, err := redisLock.TryLock("service ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tryLock)
	if !tryLock {
		return
	}
	time.Sleep(20 * time.Second)
	redisLock.UnLock()
	time.Sleep(4 * time.Second)
	fmt.Println("Done")
}
