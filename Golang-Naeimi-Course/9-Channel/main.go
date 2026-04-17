package main

import (
	"fmt"
	"time"
)

func main() {
	//SimpleUnbufferedChannel1()
	//SimpleUnbufferedChannel2()

	//BufferedChannel()
	//ForLoopInChannel()
	CheckForCloseChannelExample()
}

func SimpleUnbufferedChannel1() {
	ch := make(chan int)

	go func() {
		ch <- 22
	}()

	val := <-ch
	fmt.Println(val)
}

func SimpleUnbufferedChannel2() {
	ch := make(chan int)
	go SendDataToChannel(ch)

	for i := 0; i < 3; i++ {
		val := <-ch
		fmt.Println(val)
	}

	time.Sleep(30 * time.Second)
}

func SendDataToChannel(channel chan int) {
	println("before sending data to channel")
	channel <- 1
	println("after sending data to channel 1")
	channel <- 2
	println("before sending data to channel 2")
	channel <- 3
	println("after sending data to channel 3")
}

func BufferedChannel() {
	ch := make(chan int, 2)

	ch <- 12
	ch <- 13
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func ForLoopInChannel() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 200; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

func CheckForCloseChannelExample() {
	ch := make(chan int)
	close(ch)

	value, ok := <-ch
	fmt.Println(value, ok)

	if value, ok := <-ch; ok {
		fmt.Println("Channel Is Open %v", value)
	} else {
		fmt.Println("channel Is Close")
	}
}
