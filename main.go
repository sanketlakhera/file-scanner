package main

import (
	"fmt"
)

func main() {
	x := 42
	p := &x         // p holds the address of x
	fmt.Println(p)  // 0xc0000b6010 (some memory address)
	fmt.Println(*p) // 42
}
