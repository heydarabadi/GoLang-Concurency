package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	//NonBlockingSelect()
	//TimeOut()
	//UseContextForManageOperation()
	//MainForFanIn()
	//MainForFanOut()
	MainForTee()
}

func NonBlockingSelect() {
	ch := make(chan int)

	select {
	case value := <-ch:
		fmt.Println("Received ", value)
	default:
		fmt.Println("No value received")
	}

	fmt.Println("Close Channel ...")
}

func TimeOut() {
	ch := make(chan int)

	go func() {
		time.Sleep(5 * time.Second)
		ch <- 1
	}()

	select {
	case value := <-ch:
		fmt.Println("Received ", value)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}
}

func UseContextForManageOperation() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling Context")
		cancel()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context done", ctx.Err())
	}
}

// FanIn
func fanIn(ch1, ch2 <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		for {
			select {
			case msg := <-ch1:
				out <- msg
			case msg := <-ch2:
				out <- msg
			}
		}
	}()

	return out
}
func MainForFanIn() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- fmt.Sprintf("کانال1: پیام %d", i)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			ch2 <- fmt.Sprintf("کانال2: پیام %d", i)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	merged := fanIn(ch1, ch2)

	for i := 1; i <= 6; i++ {
		fmt.Println(<-merged)
	}
}

// FanOut
func Worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Println("Worker Started Work ", id, " - ", j)
		time.Sleep(1 * time.Second)
		fmt.Println("Worker Ended Work ", id, " - ", j)

	}
}
func MainForFanOut() {

	start := time.Now()

	jobs := make(chan int, 100)

	var wg sync.WaitGroup

	for w := 1; w <= 50; w++ {
		wg.Add(1)
		go Worker(w, jobs, &wg)
	}

	for j := 1; j <= 100; j++ {
		jobs <- j
	}
	close(jobs)
	wg.Wait()
	fmt.Println("Time Taken: ", time.Since(start))
}

// Tee
func tee(ch <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range ch {
			out1 <- val
			out2 <- val
		}
	}()
	return out1, out2
}
func MainForTee() {
	input := make(chan int)

	out1, out2 := tee(input)

	go func() {
		for val := range out1 {
			fmt.Printf("Out1 Received: %d\n", val)
		}
	}()

	go func() {
		for val := range out2 {
			fmt.Printf("Out2 Received: %d\n", val)
		}
	}()

	for i := 1; i <= 1000; i++ {
		input <- i
		time.Sleep(5 * time.Millisecond)
	}
	close(input)

	time.Sleep(15 * time.Second)

}
