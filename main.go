package main

import (
	"fmt"
	"redis-lock/lock"
	"time"
)

// git push -u origin main
func main() {
	redisLock := lock.MakeRedisLock("service")
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
	//fmt.Println(redisLock)
	//time.Sleep(4 * time.Second)
	fmt.Println("Done")
}
