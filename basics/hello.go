package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")

	// Variables
	var name string = "Gopher"
	age := 10

	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// Basic types
	var (
		isActive bool    = true
		count    int     = 42
		price    float64 = 19.99
	)

	fmt.Printf("Active: %t, Count: %d, Price: %.2f\n", isActive, count, price)
}
