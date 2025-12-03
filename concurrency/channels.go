package main

import "fmt"

func sum(numbers []int, result chan int) {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	result <- sum // Send result to channel
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6}

	resultChan := make(chan int)

	// Split work between two goroutines
	go sum(numbers[:len(numbers)/2], resultChan)
	go sum(numbers[len(numbers)/2:], resultChan)

	// Receive results from channel
	part1, part2 := <-resultChan, <-resultChan

	fmt.Printf("Part 1: %d, Part 2: %d, Total: %d\n", part1, part2, part1+part2)
}
