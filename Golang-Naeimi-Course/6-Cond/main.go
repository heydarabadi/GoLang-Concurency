package main

import (
	"fmt"
	"sync"
	"time"
)

var userList []int
var ready = false

func main() {
	cond := sync.Cond{L: &sync.Mutex{}}
	for i := 0; i < 200; i++ {
		go NewRequest(i, &cond)
	}
}

func NewRequest(userId int, cond *sync.Cond) {
	Checking(userId, cond)

	fmt.Println("NewRequest")
	cond.L.Lock()
	defer cond.L.Unlock()
	if !ready {
		cond.Wait()
	}
	fmt.Printf("User: %v  Is Ready For Streaming\n", userId)
}

func Checking(userId int, cond *sync.Cond) {
	fmt.Printf("User Id %v Waiting For Start Streaming\n", userId)
	time.Sleep(300 * time.Millisecond)
	cond.L.Lock()
	defer cond.L.Unlock()
	userList = append(userList, userId)
	if len(userList) == 66 {
		ready = true
		cond.Broadcast()
		fmt.Printf("Waiting End Stream With UserId: %v ", userId)
	}
}
