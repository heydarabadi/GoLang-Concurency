package main

import (
	"fmt"
	"sync"
)

type Person struct {
	Name string
	Age  int
}

var personPool = sync.Pool{
	New: func() interface{} {
		return &Person{}
	},
}

func main() {
	p := personPool.Get().(*Person)

	p.Name = "Ali"
	p.Age = 25
	fmt.Println(p)

	personPool.Put(p)

	p2 := personPool.Get().(*Person)

	fmt.Println(p2)

	p2.Name = ""
	p2.Age = 0

	p3 := personPool.Get().(*Person)
	fmt.Println(p3)

}
