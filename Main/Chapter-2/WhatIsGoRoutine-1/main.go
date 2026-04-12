package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		for i := 0; i < 300; i++ {
			fmt.Println("Hello World From GoRoutine")

		}
	}()

	fmt.Println("Hello World Main")
	time.Sleep(100 * time.Second)
}
