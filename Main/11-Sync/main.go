package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

}

func Mutex() {
	var mu sync.Mutex
	counter := 0

	for i := 0; i < 5; i++ {
		go func() {
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	fmt.Scanln()
	fmt.Println(counter)
}

func RWMutex() {
	var mu sync.RWMutex
	data := 0

	go func() {
		mu.Lock()
		data = 10
		mu.Unlock()
	}()

	go func() {
		mu.RLock()
		fmt.Println(data)
		mu.RUnlock()
	}()

	fmt.Scanln()
}

func WaitGroup() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			fmt.Println("worker", n)
		}(i)
	}

	wg.Wait()

	fmt.Println("done")
}

func Once() {
	var once sync.Once

	task := func() {
		fmt.Println("run only once")
	}

	for i := 0; i < 5; i++ {
		go func() {
			once.Do(task)
		}()
	}

	fmt.Scanln()
}

func Cond() {
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)

	ready := false

	go func() {
		cond.L.Lock()
		for !ready {
			cond.Wait()
		}
		fmt.Println("Goroutine awakened!")
		cond.L.Unlock()
	}()

	cond.L.Lock()
	ready = true
	cond.Signal()
	cond.L.Unlock()

	fmt.Scanln()
}

func Map() {
	var m sync.Map

	m.Store("a", 1)
	m.Store("b", 2)

	v, ok := m.Load("a")
	fmt.Println(v, ok) // 1 true

	m.Delete("b")

	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

func Atomic() {
	var counter int64 = 0

	atomic.AddInt64(&counter, 1)
	fmt.Println(counter) // 1
}
