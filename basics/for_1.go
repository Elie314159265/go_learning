/*
GoのForは少し特徴的。まずは以下のコードを見てみる。
for rangeでは for index, value := range スライス {のようになっておりindexはインデックス番号(0,1,2,3,4,・・・)のような値を返し、valueは実際の値(2,3,5)を保持する。このとき、インデックスが不要なことがあるのでこの時は以下のような書き方をする。
for _, value := range numsこれによって_はこの値は使いませんと明示することが可能。Goの思想は無駄を極限まで省くという考えなのでコードが曖昧にならないように、バクを防ぐために、コンパイル後の最適化のためにこのように使わない変数を「エラー」にする。
for _, v := range numsはGoのベストプラクティス。
*/

package main

import (
	"fmt"
)

func main() {
	nums := []int{2,3,5}
	for index, value := range nums {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}
}

