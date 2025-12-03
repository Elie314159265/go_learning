package main

import "fmt"

func main() {
	ch := make(chan int, 2) // バッファ容量 2

	fmt.Println("send 1")
	ch <- 1 // まだ空きがあるので送れる

	fmt.Println("send 2")
	ch <- 2 // ここでバッファ満杯

	fmt.Println("try to send 3 ... (this will block)")
	// ここでブロックして先へ進めない
	ch <- 3

	fmt.Println("received:", <-ch)
	fmt.Println("received:", <-ch)
	fmt.Println("received:", <-ch)
}

