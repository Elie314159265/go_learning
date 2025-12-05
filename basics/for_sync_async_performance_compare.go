package main

import (
	"fmt"
	"sync"
	"time"
)

// ダミーの重い処理,ネットワークI/Oを含む処理など
func heavyTask(n int) {
	time.Sleep(50 * time.Millisecond) // 50ms の重い仕事だと仮定
	fmt.Println("done:", n)
}

// 並列実行（非同期）
func runConcurrent(nums []int) time.Duration {
	var wg sync.WaitGroup
	start := time.Now()

	for _, v := range nums {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			heavyTask(x)
		}(v)
	}

	wg.Wait()
	return time.Since(start)
}

// 逐次実行（同期）
func runSequential(nums []int) time.Duration {
	start := time.Now()

	for _, v := range nums {
		heavyTask(v)
	}

	return time.Since(start)
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	fmt.Println("=== sequential (同期) ===")
	seq := runSequential(nums)
	fmt.Println("time:", seq)

	fmt.Println("\n=== concurrent (並行) ===")
	con := runConcurrent(nums)
	fmt.Println("time:", con)

	fmt.Printf("\nSpeedup: 約 %.1fx\n", float64(seq)/float64(con))
}

