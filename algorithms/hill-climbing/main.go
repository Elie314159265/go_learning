package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// HillClimbing は山登り法アルゴリズムを実装する構造体
type HillClimbing struct {
	currentSolution float64 // 現在の解
	currentValue    float64 // 現在の解の評価値
	stepSize        float64 // 近傍探索のステップサイズ
	maxIterations   int     // 最大反復回数
}

// NewHillClimbing はHillClimbingの新しいインスタンスを作成する
func NewHillClimbing(initialSolution, stepSize float64, maxIterations int) *HillClimbing {
	hc := &HillClimbing{
		currentSolution: initialSolution,
		stepSize:        stepSize,
		maxIterations:   maxIterations,
	}
	// 初期解の評価値を計算
	hc.currentValue = hc.evaluationFunction(initialSolution)
	return hc
}

// evaluationFunction は目的関数（最大化したい関数）
// 例: f(x) = -x^2 + 4x + 10 （放物線、最大値はx=2で14）
func (hc *HillClimbing) evaluationFunction(x float64) float64 {
	return -math.Pow(x, 2) + 4*x + 10
}

// generateNeighbor は現在の解の近傍解を生成する
// ランダムに±stepSize の範囲で新しい解を生成
func (hc *HillClimbing) generateNeighbor() float64 {
	// -stepSize から +stepSize の範囲でランダムな値を生成
	delta := (rand.Float64()*2 - 1) * hc.stepSize
	return hc.currentSolution + delta
}

// Solve は山登り法を実行して最適解を探索する
func (hc *HillClimbing) Solve() {
	fmt.Println("=== Hill Climbing アルゴリズム開始 ===")
	fmt.Printf("初期解: x = %.4f, f(x) = %.4f\n\n", hc.currentSolution, hc.currentValue)

	improvementCount := 0 // 改善された回数

	// 最大反復回数まで探索を繰り返す
	for iteration := 0; iteration < hc.maxIterations; iteration++ {
		// 1. 近傍解を生成
		neighbor := hc.generateNeighbor()

		// 2. 近傍解の評価値を計算
		neighborValue := hc.evaluationFunction(neighbor)

		// 3. 近傍解が現在の解より良い場合、更新する
		if neighborValue > hc.currentValue {
			fmt.Printf("反復 %d: 改善発見!\n", iteration+1)
			fmt.Printf("  現在解: x = %.4f, f(x) = %.4f\n", hc.currentSolution, hc.currentValue)
			fmt.Printf("  新解:   x = %.4f, f(x) = %.4f\n", neighbor, neighborValue)
			fmt.Printf("  改善量: %.4f\n\n", neighborValue-hc.currentValue)

			// 解を更新
			hc.currentSolution = neighbor
			hc.currentValue = neighborValue
			improvementCount++
		}
	}

	// 結果を表示
	fmt.Println("=== Hill Climbing アルゴリズム終了 ===")
	fmt.Printf("最終解: x = %.4f, f(x) = %.4f\n", hc.currentSolution, hc.currentValue)
	fmt.Printf("総改善回数: %d / %d 回の反復\n", improvementCount, hc.maxIterations)

	// 理論的な最適解と比較
	optimalX := 2.0
	optimalValue := hc.evaluationFunction(optimalX)
	fmt.Printf("\n理論的最適解: x = %.4f, f(x) = %.4f\n", optimalX, optimalValue)
	fmt.Printf("誤差: %.4f\n", math.Abs(hc.currentSolution-optimalX))
}

// SolveWithRandomRestart はランダムリスタート付き山登り法を実行
// 複数の初期解から開始し、最良の解を見つける
func SolveWithRandomRestart(numRestarts int, searchRange float64, stepSize float64, maxIterations int) {
	fmt.Println("\n=== ランダムリスタート Hill Climbing ===")
	fmt.Printf("リスタート回数: %d\n\n", numRestarts)

	var bestSolution float64
	var bestValue float64 = math.Inf(-1) // 負の無限大で初期化

	// 複数の初期解から探索を実行
	for i := 0; i < numRestarts; i++ {
		// ランダムな初期解を生成（-searchRange ~ +searchRange）
		initialSolution := (rand.Float64()*2 - 1) * searchRange

		fmt.Printf("--- リスタート %d ---\n", i+1)
		hc := NewHillClimbing(initialSolution, stepSize, maxIterations)

		// 簡易版の探索（詳細出力なし）
		for iteration := 0; iteration < maxIterations; iteration++ {
			neighbor := hc.generateNeighbor()
			neighborValue := hc.evaluationFunction(neighbor)

			if neighborValue > hc.currentValue {
				hc.currentSolution = neighbor
				hc.currentValue = neighborValue
			}
		}

		fmt.Printf("結果: x = %.4f, f(x) = %.4f\n\n", hc.currentSolution, hc.currentValue)

		// 最良解を更新
		if hc.currentValue > bestValue {
			bestSolution = hc.currentSolution
			bestValue = hc.currentValue
		}
	}

	fmt.Println("=== 最終結果 ===")
	fmt.Printf("最良解: x = %.4f, f(x) = %.4f\n", bestSolution, bestValue)
}

func main() {
	// 乱数シードを初期化
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Hill Climbing アルゴリズムのデモ")
	fmt.Println("目的関数: f(x) = -x^2 + 4x + 10")
	fmt.Println("目標: f(x) の最大値を見つける\n")

	// パラメータ設定
	initialSolution := 0.0  // 初期解
	stepSize := 0.5         // 近傍探索のステップサイズ
	maxIterations := 100    // 最大反復回数

	// 基本的な Hill Climbing を実行
	hc := NewHillClimbing(initialSolution, stepSize, maxIterations)
	hc.Solve()

	// ランダムリスタート版を実行
	SolveWithRandomRestart(5, 10.0, 0.5, 50)
}
