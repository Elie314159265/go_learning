package main

import "fmt"

// Fibonacci 数を送り続けるが、quit が来たら終了する
func fibonacci(out chan int, quit chan struct{}) {
	x, y := 0, 1
	for {
		select {
		case out <- x: // out に値を送れるなら送る
			x, y = y, x+y
		case <-quit: // quit チャネルに何か来たら終了
			fmt.Println("stop signal received")
			close(out)
			return
		}
	}
}

func main() {
	out := make(chan int)
	quit := make(chan struct{})

	// out から10個受け取って quit の合図を送る
	go func() {
		for i := 0; i < 30; i++ {
			fmt.Println(<-out)
		}
		// struct{}はサイズ０バイトなので最も軽量な通知チャネルとして利用される
		quit <- struct{}{}
	}()

	fibonacci(out, quit)
}
