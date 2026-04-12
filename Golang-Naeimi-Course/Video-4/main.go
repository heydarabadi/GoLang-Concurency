package main

import (
	"fmt"
	"strconv"
	"sync"
)

var todoList []string

func main() {
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		urlNew := "http://google.com" + "/" + strconv.Itoa(i) + "/"
		wg.Add(1)
		go GetUrl(urlNew, &wg)
	}

	wg.Wait()

	fmt.Printf("%v \n", todoList)
}

func GetUrl(url string, wg *sync.WaitGroup) {
	todoList = append(todoList, url)
	defer wg.Done()
}
