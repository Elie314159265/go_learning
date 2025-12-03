package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values from the tree to channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walkFn func(t *tree.Tree)
	walkFn = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walkFn(t.Left)     // 左
		ch <- t.Value      // 自分
		walkFn(t.Right)    // 右（中順巡回：in-order traversal）
	}
	walkFn(t)
	close(ch)
}

// Same checks whether t1 and t2 store the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if ok1 != ok2 {
			return false // 要素数が違う
		}
		if !ok1 { // どちらも読み切った
			return true
		}
		if v1 != v2 {
			return false // 値が違う
		}
	}
}

func main() {
	// Walk のテスト
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for v := range ch {
		fmt.Println(v)
	}

	// Same のテスト
	fmt.Println("Same(1, 1) =", Same(tree.New(1), tree.New(1)))
	fmt.Println("Same(1, 2) =", Same(tree.New(1), tree.New(2)))
}
/*
① Walk(t, ch) は「中順巡回（in-order traversal）」
中順巡回は 左 → 自分 → 右 という順。
二分探索木（BST）でこれをやると 必ずソートされた順に値が出てくる。

② チャネルは最後に close(ch) する
受け取る側（range ch）が安全に終了できるようにします。
Close を忘れると goroutine がデッドロックになります。

③ Same 関数では「2つのチャネルを同時に読み比べる」
他言語だと木構造を比較するのはやっかいですが、Go は Walk → チャネル → 比較 でシンプルになる。値が違えば false、読み終わり状態（ok1 / ok2）が違えばfalse、全部同じなら true
チャネルの “逐次比較” は Go らしい発想。

           Tree
            4
          /   \
         2     6
        / \   / \
       1  3  5   7
中順巡回すると：
1 → 2 → 3 → 4 → 5 → 6 → 7

これをチャネルに流すだけ。
Same はこのストリームを 2本同時に比較するだけ

これの何が良い？
- 「木構造 → ストリーム変換」によって「比較処理」が簡素化できる
- 並行処理(Goroutine/Channel)を自然に使える
- データサイズが大きくても、値をストリーム比較するので効率が良い（全体をメモリに持たない）
*/
