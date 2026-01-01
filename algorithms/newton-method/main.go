package main

import (
	"fmt"
	"math"
)

// NewtonMethod はニュートン法を実装する構造体
type NewtonMethod struct {
	initialGuess  float64 // 初期推定値
	tolerance     float64 // 許容誤差
	maxIterations int     // 最大反復回数
}

// NewNewtonMethod はNewtonMethodの新しいインスタンスを作成する
func NewNewtonMethod(initialGuess, tolerance float64, maxIterations int) *NewtonMethod {
	return &NewtonMethod{
		initialGuess:  initialGuess,
		tolerance:     tolerance,
		maxIterations: maxIterations,
	}
}

// function は求根したい関数 f(x)
// 例1: f(x) = x^2 - 2 （√2 を求める）
func (nm *NewtonMethod) function(x float64) float64 {
	return x*x - 2
}

// derivative は関数の導関数 f'(x)
// f(x) = x^2 - 2 の導関数は f'(x) = 2x
func (nm *NewtonMethod) derivative(x float64) float64 {
	return 2 * x
}

// Solve はニュートン法を実行して方程式の解を見つける
// ニュートン法の公式: x_{n+1} = x_n - f(x_n) / f'(x_n)
func (nm *NewtonMethod) Solve() (float64, int, bool) {
	fmt.Println("=== ニュートン法アルゴリズム開始 ===")
	fmt.Println("方程式: f(x) = x² - 2 = 0")
	fmt.Println("目標: √2 の値を求める")
	fmt.Printf("初期推定値: x₀ = %.6f\n", nm.initialGuess)
	fmt.Printf("許容誤差: %.10f\n\n", nm.tolerance)

	x := nm.initialGuess // 現在の推定値

	// 反復計算
	for iteration := 0; iteration < nm.maxIterations; iteration++ {
		// 1. 現在の x における関数値 f(x) を計算
		fx := nm.function(x)

		// 2. 現在の x における導関数 f'(x) を計算
		fpx := nm.derivative(x)

		// 3. 導関数が0に近い場合（接線が水平）、計算を中止
		if math.Abs(fpx) < 1e-10 {
			fmt.Printf("警告: 導関数が0に近いため計算を中止します（x = %.6f）\n", x)
			return x, iteration, false
		}

		// 4. ニュートン法の公式を適用して次の推定値を計算
		// x_{n+1} = x_n - f(x_n) / f'(x_n)
		xNext := x - fx/fpx

		// 進捗を表示
		fmt.Printf("反復 %d:\n", iteration+1)
		fmt.Printf("  x_%d     = %.10f\n", iteration, x)
		fmt.Printf("  f(x_%d)  = %.10f\n", iteration, fx)
		fmt.Printf("  f'(x_%d) = %.10f\n", iteration, fpx)
		fmt.Printf("  x_%d     = %.10f\n", iteration+1, xNext)
		fmt.Printf("  誤差     = %.10f\n\n", math.Abs(xNext-x))

		// 5. 収束判定：前回と今回の推定値の差が許容誤差以下なら終了
		if math.Abs(xNext-x) < nm.tolerance {
			fmt.Println("=== 収束成功 ===")
			fmt.Printf("最終解: x = %.10f\n", xNext)
			fmt.Printf("検算: f(%.10f) = %.15f\n", xNext, nm.function(xNext))
			fmt.Printf("反復回数: %d\n", iteration+1)
			return xNext, iteration + 1, true
		}

		// 6. 次の反復へ
		x = xNext
	}

	// 最大反復回数に達しても収束しなかった場合
	fmt.Println("=== 収束失敗 ===")
	fmt.Printf("最大反復回数 %d に達しました\n", nm.maxIterations)
	fmt.Printf("最終推定値: x = %.10f\n", x)
	return x, nm.maxIterations, false
}

// Example2_CubicEquation は三次方程式の解を求める例
// f(x) = x^3 - x - 2 = 0 の解を求める（解の一つは x ≈ 1.521）
type CubicEquation struct {
	initialGuess  float64
	tolerance     float64
	maxIterations int
}

// NewCubicEquation はCubicEquationの新しいインスタンスを作成
func NewCubicEquation(initialGuess, tolerance float64, maxIterations int) *CubicEquation {
	return &CubicEquation{
		initialGuess:  initialGuess,
		tolerance:     tolerance,
		maxIterations: maxIterations,
	}
}

// function は f(x) = x^3 - x - 2
func (ce *CubicEquation) function(x float64) float64 {
	return x*x*x - x - 2
}

// derivative は f'(x) = 3x^2 - 1
func (ce *CubicEquation) derivative(x float64) float64 {
	return 3*x*x - 1
}

// Solve はニュートン法で三次方程式を解く
func (ce *CubicEquation) Solve() (float64, int, bool) {
	fmt.Println("\n=== 三次方程式の求解 ===")
	fmt.Println("方程式: f(x) = x³ - x - 2 = 0")
	fmt.Printf("初期推定値: x₀ = %.6f\n\n", ce.initialGuess)

	x := ce.initialGuess

	for iteration := 0; iteration < ce.maxIterations; iteration++ {
		fx := ce.function(x)
		fpx := ce.derivative(x)

		if math.Abs(fpx) < 1e-10 {
			fmt.Printf("警告: 導関数が0に近いため計算を中止します\n")
			return x, iteration, false
		}

		xNext := x - fx/fpx

		fmt.Printf("反復 %d: x = %.10f, f(x) = %.10f\n", iteration+1, xNext, ce.function(xNext))

		if math.Abs(xNext-x) < ce.tolerance {
			fmt.Println("\n=== 収束成功 ===")
			fmt.Printf("最終解: x = %.10f\n", xNext)
			fmt.Printf("検算: f(%.10f) = %.15f\n", xNext, ce.function(xNext))
			fmt.Printf("反復回数: %d\n", iteration+1)
			return xNext, iteration + 1, true
		}

		x = xNext
	}

	fmt.Println("収束失敗")
	return x, ce.maxIterations, false
}

// SquareRoot はニュートン法を使って平方根を計算する関数
// a の平方根を求める（x^2 - a = 0 の解）
func SquareRoot(a, initialGuess, tolerance float64, maxIterations int) float64 {
	fmt.Printf("\n=== √%.0f の計算 ===\n", a)
	fmt.Printf("初期推定値: %.6f\n\n", initialGuess)

	x := initialGuess

	for iteration := 0; iteration < maxIterations; iteration++ {
		// f(x) = x^2 - a
		// f'(x) = 2x
		// x_{n+1} = x_n - (x_n^2 - a) / (2 * x_n)
		//         = (x_n + a/x_n) / 2  （簡略化した形）

		xNext := (x + a/x) / 2

		fmt.Printf("反復 %d: x = %.10f, 誤差 = %.10f\n",
			iteration+1, xNext, math.Abs(xNext-x))

		if math.Abs(xNext-x) < tolerance {
			fmt.Printf("\n結果: √%.0f = %.10f\n", a, xNext)
			fmt.Printf("検算: %.10f² = %.10f\n", xNext, xNext*xNext)
			fmt.Printf("math.Sqrt(%.0f) = %.10f\n", a, math.Sqrt(a))
			fmt.Printf("誤差: %.15f\n", math.Abs(xNext-math.Sqrt(a)))
			return xNext
		}

		x = xNext
	}

	return x
}

// VisualizeTangentLine は接線の概念を視覚的に説明する
func VisualizeTangentLine() {
	fmt.Println("\n=== ニュートン法の幾何学的解釈 ===")
	fmt.Println("ニュートン法は、関数の接線を使って根に近づく方法です")
	fmt.Println()
	fmt.Println("  f(x) |")
	fmt.Println("       |")
	fmt.Println("     + |      ●  <- (x_n, f(x_n))")
	fmt.Println("       |     /|")
	fmt.Println("       |    / |")
	fmt.Println("       |   /  |  接線")
	fmt.Println("       |  /   |")
	fmt.Println("   0 --|-----●---------> x")
	fmt.Println("       |    x_n+1  x_n")
	fmt.Println()
	fmt.Println("1. 点 (x_n, f(x_n)) における接線を引く")
	fmt.Println("2. 接線とx軸の交点が x_n+1 となる")
	fmt.Println("3. x_n+1 を新しい推定値として、1に戻る")
	fmt.Println()
	fmt.Println("公式: x_{n+1} = x_n - f(x_n) / f'(x_n)")
	fmt.Println("      └─────┘   └──┘   └────┘")
	fmt.Println("       次の値   現在値  接線の傾き")
}

func main() {
	fmt.Println("ニュートン法（Newton's Method）のデモ")
	fmt.Println("==========================================\n")

	// 幾何学的解釈の説明
	VisualizeTangentLine()

	// 例1: √2 を求める (x^2 - 2 = 0)
	fmt.Println("\n【例1】平方根の計算")
	nm := NewNewtonMethod(1.0, 1e-10, 20)
	solution, _, converged := nm.Solve()

	if converged {
		fmt.Printf("\n理論値との比較:\n")
		fmt.Printf("計算結果:  %.10f\n", solution)
		fmt.Printf("math.Sqrt: %.10f\n", math.Sqrt(2))
		fmt.Printf("誤差:      %.15f\n", math.Abs(solution-math.Sqrt(2)))
	}

	// 例2: 簡略化した平方根計算
	SquareRoot(10, 3.0, 1e-10, 20)

	// 例3: 三次方程式の解
	fmt.Println("\n【例2】三次方程式の求解")
	cubic := NewCubicEquation(2.0, 1e-10, 20)
	cubic.Solve()

	// 収束失敗の例
	fmt.Println("\n【例3】収束しない場合の例")
	fmt.Println("初期推定値が不適切な場合、収束しないことがあります")
	badNm := NewNewtonMethod(0.0, 1e-10, 5)
	_, _, converged = badNm.Solve()
	if !converged {
		fmt.Println("適切な初期推定値を選ぶことが重要です！")
	}
}
