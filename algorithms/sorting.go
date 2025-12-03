package main

import "fmt"

// BubbleSort implements bubble sort algorithm
func BubbleSort(arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	copy(result, arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// QuickSort implements quick sort algorithm
func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	result := make([]int, len(arr))
	copy(result, arr)

	pivot := result[len(result)/2]
	var left, right, middle []int

	for _, v := range result {
		if v < pivot {
			left = append(left, v)
		} else if v > pivot {
			right = append(right, v)
		} else {
			middle = append(middle, v)
		}
	}

	return append(append(QuickSort(left), middle...), QuickSort(right)...)
}

func main() {
	arr := []int{64, 34, 25, 12, 22, 11, 90}

	fmt.Println("Original:", arr)
	fmt.Println("Bubble Sort:", BubbleSort(arr))
	fmt.Println("Quick Sort:", QuickSort(arr))
}
