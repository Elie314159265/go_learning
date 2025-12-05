/*
sync.WaitGroupのAdd,Done,Waitについて

wg.Add(n): 待つべき作業を増やす。引数は正の値で通常１。
wg.Done: 作業の完了を通知する。(wg.Add(-1)と同義)
wg.Wait(): カウントが0になるまでブロックする。全作業が終わるまで待つ。

注意点
１．ゴルーチンを起動する前にwg.Add(1)を呼ぶ。起動後にAddするとゴルーチンがすぐにDoneを呼んでAddより先にカウントが０→負になるとPanicになるリスクがある。したがってwg.Add → go func()の流れ

２．ゴルーチン内ではdefer wg.Done()を使う。早期リターンやパニック時でもカウントが確実に減る。

３．wg.Add()の呼び出しはwg.Wait()が既に返ったあとに行わない。Wait()が返ったあと(全タスク終了)にAddすると扱いが混乱する。一般にWaitのあとにAddの設計は避けるべき。

４．
*/
package main

import (
	"fmt"
	"sync"
)

// WaitGroupは状態を共有する必要があるためポインタで渡す。
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("worker %d started\n", id)
	// 仕事をする
	fmt.Printf("worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup

	const n = 3
	var i int
	for i = 1; i <= n; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("all workers finished")
}
