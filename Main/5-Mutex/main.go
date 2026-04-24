package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	UseRwMutex()
}

// Simple Example
func example1() {
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	wg.Add(1000000)

	counter := 0
	for i := 0; i < 1000000; i++ {
		go func() {
			defer wg.Done()
			mx.Lock()
			counter++
			defer mx.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("%v \n", counter)
}

type Employee struct {
	name   string
	salary int64
}

func example2() {
	var resourceMoney int64 = 12000000
	wg := sync.WaitGroup{}
	wg.Add(500)

	mx := sync.Mutex{}

	employeeList := []Employee{}
	for i := 0; i < 500; i++ {
		rand.Seed(time.Now().UnixNano())

		employeeList = append(employeeList, Employee{"hasan-" + strconv.Itoa(i), 50000})
	}

	for _, employee := range employeeList {

		fmt.Println(employee.name)
		go func(employee Employee) {
			defer wg.Done()
			mx.Lock()
			if employee.salary < resourceMoney {

				resourceMoney -= employee.salary
				mx.Unlock()
			}

		}(employee)
	}

	wg.Wait()

	println(resourceMoney)
}

func example3() {
	var resourceMoney int64 = 12000000
	wg := sync.WaitGroup{}
	wg.Add(500)

	employeeList := []Employee{}
	for i := 0; i < 500; i++ {
		rand.Seed(time.Now().UnixNano())

		employeeList = append(employeeList, Employee{"hasan-" + strconv.Itoa(i), 50000})
	}

	for _, employee := range employeeList {

		fmt.Println(employee.name)
		go func(employee Employee) {
			defer wg.Done()
			if employee.salary < resourceMoney {
				atomic.AddInt64(&resourceMoney, employee.salary)
			}

		}(employee)
	}

	wg.Wait()

	println(resourceMoney)
}

func CreateDeadLock() {
	var mu sync.Mutex

	fmt.Println("Lock #1")
	mu.Lock()

	fmt.Println("Lock #2 (DeadLock Happens!)")
	mu.Lock()

	fmt.Println("This Will Never Be Printed")
}

func ResolveDeadLock() {
	m1 := sync.Mutex{}
	m2 := sync.Mutex{}

	go func() {
		m1.Lock()
		fmt.Println("G1 Got m1")
		time.Sleep(50 * time.Millisecond)

		m2.Lock()
		fmt.Println("G2 Got m2")

		m1.Unlock()
		m2.Unlock()
	}()

	go func() {
		m1.Lock()
		fmt.Println("G2 Got m1")

		time.Sleep(50 * time.Millisecond)

		m2.Lock()
		fmt.Println("G2 Got m2")

		m2.Unlock()
		m1.Unlock()
	}()

	time.Sleep(1 * time.Second)

}

func UseSafeUnlockAndDoneWithDefer() {
	m1 := sync.Mutex{}
	counter := 0
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		m1.Lock()
		defer m1.Unlock()
		counter++
		counter += 10
		fmt.Println("Counter Updated: ", counter)
	}()

	wg.Wait()
}

func UseRwMutex() {
	rwMutext := sync.RWMutex{}
	value := 100
	wg := sync.WaitGroup{}
	wg.Add(11)

	// 10 reader
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			rwMutext.RLock()
			fmt.Printf("value is : %v with Id: %v \n", value, id)
			rwMutext.RUnlock()
		}(i)

	}

	// 1 writer
	go func() {
		defer wg.Done()
		rwMutext.Lock()
		value = 10
		fmt.Println("Writer updated value: ", value)
		rwMutext.Unlock()
	}()

	wg.Wait()

}
