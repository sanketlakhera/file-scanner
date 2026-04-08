package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("hello from goroutine")
}

func main() {
	go sayHello() // we can mark any programm as goroutine by adding "go" keyword before the function call
	fmt.Println("hello from main")
	time.Sleep(time.Millisecond) // problem is that main function will exit before the goroutine function is called, that's why this little hack
}
