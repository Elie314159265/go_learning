package main

import (
	"fmt"
	"time"
)

func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Hello, %s! (%d)\n", name, i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Goroutines - concurrent execution
	go sayHello("Alice")
	go sayHello("Bob")

	// Wait for goroutines to finish
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Done!")
}
