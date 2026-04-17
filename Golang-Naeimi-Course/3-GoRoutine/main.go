package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	go task1()
	go task2()
	go task3()

	value := 0
	go func() {
		value++
	}()
	go func() {
		value += 2
	}()
	go func() {
		value += 3
	}()
	go func() {
		value += 3
	}()

	fmt.Println(value)

	fmt.Println(runtime.NumGoroutine())

	time.Sleep(time.Second)
}

func task1() {
	println("Task 1")
}

func task2() {
	println("Task 2")
}

func task3() {
	println("Task 3")
}
