package main

import (
	"fmt"
	"sync"
	"time"
)

// =========================
// ダミーの I/O（超重め）
// =========================
func fakeIO(id int) {
	// ネットワークやディスク待ちを模倣（100ms）
	time.Sleep(100 * time.Millisecond)
}

// =========================
// 同期処理
// =========================
func runSync(n int) time.Duration {
	start := time.Now()
	for i := 0; i < n; i++ {
		fakeIO(i)
	}
	return time.Since(start)
}

// =========================
// 非同期処理
// =========================
func runAsync(n int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			fakeIO(id)
		}(i)
	}

	wg.Wait()
	return time.Since(start)
}

func main() {
	N := 100 // 1000件のI/Oタスク

	fmt.Println("=== 同期処理 ===")
	syncTime := runSync(N)
	fmt.Println("time:", syncTime)

	fmt.Println("\n=== 非同期処理 ===")
	asyncTime := runAsync(N)
	fmt.Println("time:", asyncTime)

	fmt.Printf("\nSpeedup: 約 %.1fx\n", float64(syncTime)/float64(asyncTime))
}

