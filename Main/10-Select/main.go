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
	//MainForTee()
	//MainForBridge()
	//MainForMultiplex()
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

// Bridge

func bridge(channels <-chan <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for ch := range channels {
			for v := range ch {
				out <- v
			}
		}
	}()

	return out
}

func MainForBridge() {
	channels := make(chan (<-chan int))

	go func() {
		defer close(channels)

		a := make(chan int)
		go func() {
			defer close(a)
			a <- 1
			a <- 2
		}()

		b := make(chan int)
		go func() {
			defer close(b)
			b <- 3
			b <- 4
		}()

		channels <- a
		channels <- b
	}()

	for v := range bridge(channels) {
		fmt.Println(v)
	}
}

// Multiplex

func multiplex(a, b <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v

			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()

	return out
}

func MainForMultiplex() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)
		a <- 1
		a <- 2
	}()

	go func() {
		defer close(b)
		b <- 3
		b <- 4
	}()

	for v := range multiplex(a, b) {
		fmt.Println(v)
	}
}

// مدیریت Backpressure با Buffer

func Backpressure() {
	ch := make(chan int, 2)
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Println("Send:", i)
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		time.Sleep(time.Second * 1)
		fmt.Println("Recv:", v)
	}
}
