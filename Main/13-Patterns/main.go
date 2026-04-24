package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

}

type Task struct {
	ID int
}

// Worker Pool

func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}
			// کار سنگین
			time.Sleep(500 * time.Millisecond)
			results <- t.ID * 2
		}
	}
}

func MainForWorkerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan Task, 10) // Work Queue
	results := make(chan int, 10)

	// 3 Worker ثابت
	for w := 1; w <= 3; w++ {
		go worker(ctx, w, tasks, results)
	}

	// ارسال 10 وظیفه
	go func() {
		for i := 1; i <= 10; i++ {
			tasks <- Task{ID: i}
		}
		close(tasks)
	}()

	// مصرف نتایج
	for i := 0; i < 10; i++ {
		fmt.Println(<-results)
	}
}

////////////////////

// Dynamic Worker Pool

func dynamicWorker(ctx context.Context, id int, tasks <-chan Task, results chan<- int, stop chan<- int) {
	for {
		select {
		case <-ctx.Done():
			stop <- id
			return
		case t, ok := <-tasks:
			if !ok {
				stop <- id
				return
			}
			time.Sleep(300 * time.Millisecond)
			results <- t.ID * 2
		}
	}
}

func manager(ctx context.Context, tasks chan Task, results chan int) {
	minWorkers := 2
	maxWorkers := 8
	queueThreshold := 5

	activeWorkers := 0
	stop := make(chan int)

	spawn := func() {
		activeWorkers++
		go dynamicWorker(ctx, activeWorkers, tasks, results, stop)
	}

	for i := 0; i < minWorkers; i++ {
		spawn()
	}

	for {
		select {
		case <-ctx.Done():
			return

		default:
			l := len(tasks)
			if l > queueThreshold && activeWorkers < maxWorkers {
				spawn()
			}

			if l == 0 && activeWorkers > minWorkers {
				cancelFunc := ctx.Done()
				_ = cancelFunc
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

////////////////////

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

////////////////////

//Bounded Parallelism

func doWork(i int) {
	time.Sleep(1 * time.Second)
	fmt.Println("done", i)
}

func BoundedParallelism() {
	sem := make(chan struct{}, 5)
	for i := 0; i < 20; i++ {
		sem <- struct{}{}

		go func(i int) {
			defer func() { <-sem }()
			doWork(i)
		}(i)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}

////////////////////

// Token Bucket
func TokenBucket(rate int, burst int) chan struct{} {
	bucket := make(chan struct{}, burst)

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()

		for range ticker.C {
			select {
			case bucket <- struct{}{}:
			default:
			}
		}
	}()

	return bucket
}

func MainForTokenBucket() {
	bucket := TokenBucket(5, 10)

	for i := 0; i < 20; i++ {
		<-bucket
		fmt.Println("request:", i)
	}

}

////////////////////
