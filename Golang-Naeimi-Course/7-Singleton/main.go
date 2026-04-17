package main

import (
	"fmt"
	"sync"
)

type Config struct {
	ConnectionString string
}

var (
	mx     = sync.Mutex{}
	config *Config
	Once   sync.Once
)

func GetConfigWithMutex() *Config {
	mx.Lock()
	defer mx.Unlock()
	if config == nil {
		config = &Config{ConnectionString: "teset"}
	}
	return config
}

func GetConfigWithOnce() *Config {
	Once.Do(func() {
		config = &Config{ConnectionString: "teset"}
	})
	return config
}

func main() {
	for i := 0; i < 100; i++ {

		//  With Lock
		con := GetConfigWithMutex()

		// With Once
		//con := GetConfigWithOnce()

		fmt.Printf("Pointer In Ram Is : %p \n", con)
	}
}
