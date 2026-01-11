package main

import (
	"fmt"
	"math"
)

// DataPoint はデータ点 (x, y) を表す構造体
type DataPoint struct {
	X float64
	Y float64
}

// LinearRegression は最小二乗法による1次式近似を実装する構造体
type LinearRegression struct {
	data   []DataPoint // データ点の配列
	a      float64     // 傾き
	b      float64     // 切片
	fitted bool        // フィッティング済みかどうか
}

// NewLinearRegression はLinearRegressionの新しいインスタンスを作成する
func NewLinearRegression(data []DataPoint) *LinearRegression {
	return &LinearRegression{
		data:   data,
		fitted: false,
	}
}

// Fit は最小二乗法で1次式 y = ax + b のパラメータを求める
// 正規方程式を使用:
//   a = (n·Σxy - Σx·Σy) / (n·Σx² - (Σx)²)
//   b = (Σy - a·Σx) / n
func (lr *LinearRegression) Fit() {
	fmt.Println("=== 1次式近似（線形回帰）開始 ===")
	fmt.Println("目標: y = ax + b の形で近似")
	fmt.Printf("データ点数: %d\n\n", len(lr.data))

	// データ点を表示
	fmt.Println("入力データ:")
	for i, point := range lr.data {
		fmt.Printf("  点%d: (x=%.2f, y=%.2f)\n", i+1, point.X, point.Y)
	}
	fmt.Println()

	n := float64(len(lr.data))

	// 各種和を計算
	var sumX, sumY, sumXY, sumX2 float64
	for _, point := range lr.data {
		sumX += point.X
		sumY += point.Y
		sumXY += point.X * point.Y
		sumX2 += point.X * point.X
	}

	// 計算過程を表示
	fmt.Println("中間計算:")
	fmt.Printf("  n     = %.0f\n", n)
	fmt.Printf("  Σx    = %.4f\n", sumX)
	fmt.Printf("  Σy    = %.4f\n", sumY)
	fmt.Printf("  Σxy   = %.4f\n", sumXY)
	fmt.Printf("  Σx²   = %.4f\n", sumX2)
	fmt.Println()

	// 正規方程式でパラメータを計算
	// a = (n·Σxy - Σx·Σy) / (n·Σx² - (Σx)²)
	numerator := n*sumXY - sumX*sumY
	denominator := n*sumX2 - sumX*sumX

	if math.Abs(denominator) < 1e-10 {
		fmt.Println("警告: 分母が0に近いため計算できません")
		return
	}

	lr.a = numerator / denominator
	lr.b = (sumY - lr.a*sumX) / n
	lr.fitted = true

	// 結果を表示
	fmt.Println("=== 計算結果 ===")
	fmt.Printf("近似式: y = %.4fx + %.4f\n", lr.a, lr.b)
	fmt.Println()

	// 各点での予測値と誤差を計算
	fmt.Println("予測値と誤差:")
	var sumSquaredError float64
	for i, point := range lr.data {
		predicted := lr.Predict(point.X)
		error := point.Y - predicted
		sumSquaredError += error * error
		fmt.Printf("  点%d: x=%.2f, 実測値=%.2f, 予測値=%.4f, 誤差=%.4f\n",
			i+1, point.X, point.Y, predicted, error)
	}

	// 決定係数 R² を計算
	meanY := sumY / n
	var totalSumSquares float64
	for _, point := range lr.data {
		totalSumSquares += (point.Y - meanY) * (point.Y - meanY)
	}

	r2 := 1.0 - (sumSquaredError / totalSumSquares)

	fmt.Println()
	fmt.Printf("二乗誤差の和: %.6f\n", sumSquaredError)
	fmt.Printf("決定係数 R²: %.6f\n", r2)
	fmt.Println("(R²が1に近いほど、データへの当てはまりが良い)")
	fmt.Println()
}

// Predict は与えられたxに対してy値を予測する
func (lr *LinearRegression) Predict(x float64) float64 {
	if !lr.fitted {
		fmt.Println("警告: Fit()を先に実行してください")
		return 0
	}
	return lr.a*x + lr.b
}

// QuadraticRegression は最小二乗法による2次式近似を実装する構造体
type QuadraticRegression struct {
	data   []DataPoint // データ点の配列
	a      float64     // x²の係数
	b      float64     // xの係数
	c      float64     // 定数項
	fitted bool        // フィッティング済みかどうか
}

// NewQuadraticRegression はQuadraticRegressionの新しいインスタンスを作成する
func NewQuadraticRegression(data []DataPoint) *QuadraticRegression {
	return &QuadraticRegression{
		data:   data,
		fitted: false,
	}
}

// Fit は最小二乗法で2次式 y = ax² + bx + c のパラメータを求める
// 正規方程式（連立方程式）を解く:
//   Σy    = a·Σx² + b·Σx  + c·n
//   Σxy   = a·Σx³ + b·Σx² + c·Σx
//   Σx²y  = a·Σx⁴ + b·Σx³ + c·Σx²
func (qr *QuadraticRegression) Fit() {
	fmt.Println("\n=== 2次式近似（放物線回帰）開始 ===")
	fmt.Println("目標: y = ax² + bx + c の形で近似")
	fmt.Printf("データ点数: %d\n\n", len(qr.data))

	// データ点を表示
	fmt.Println("入力データ:")
	for i, point := range qr.data {
		fmt.Printf("  点%d: (x=%.2f, y=%.2f)\n", i+1, point.X, point.Y)
	}
	fmt.Println()

	n := float64(len(qr.data))

	// 各種和を計算
	var sumX, sumY, sumX2, sumX3, sumX4, sumXY, sumX2Y float64
	for _, point := range qr.data {
		x := point.X
		y := point.Y
		x2 := x * x
		x3 := x2 * x
		x4 := x2 * x2

		sumX += x
		sumY += y
		sumX2 += x2
		sumX3 += x3
		sumX4 += x4
		sumXY += x * y
		sumX2Y += x2 * y
	}

	// 計算過程を表示
	fmt.Println("中間計算:")
	fmt.Printf("  n     = %.0f\n", n)
	fmt.Printf("  Σx    = %.4f\n", sumX)
	fmt.Printf("  Σy    = %.4f\n", sumY)
	fmt.Printf("  Σx²   = %.4f\n", sumX2)
	fmt.Printf("  Σx³   = %.4f\n", sumX3)
	fmt.Printf("  Σx⁴   = %.4f\n", sumX4)
	fmt.Printf("  Σxy   = %.4f\n", sumXY)
	fmt.Printf("  Σx²y  = %.4f\n", sumX2Y)
	fmt.Println()

	// 正規方程式の係数行列を構築
	// | Σx⁴  Σx³  Σx² | | a |   | Σx²y |
	// | Σx³  Σx²  Σx  | | b | = | Σxy  |
	// | Σx²  Σx   n   | | c |   | Σy   |

	// クラメルの公式で解く
	// 係数行列の行列式
	det := sumX4*(sumX2*n-sumX*sumX) -
		sumX3*(sumX3*n-sumX*sumX2) +
		sumX2*(sumX3*sumX-sumX2*sumX2)

	if math.Abs(det) < 1e-10 {
		fmt.Println("警告: 行列式が0に近いため計算できません")
		return
	}

	// クラメルの公式で a, b, c を計算
	detA := sumX2Y*(sumX2*n-sumX*sumX) -
		sumXY*(sumX3*n-sumX*sumX2) +
		sumY*(sumX3*sumX-sumX2*sumX2)

	detB := sumX4*(sumXY*n-sumY*sumX) -
		sumX3*(sumX2Y*n-sumY*sumX2) +
		sumX2*(sumX2Y*sumX-sumXY*sumX2)

	detC := sumX4*(sumX2*sumY-sumX*sumXY) -
		sumX3*(sumX3*sumY-sumX*sumX2Y) +
		sumX2*(sumX3*sumXY-sumX2*sumX2Y)

	qr.a = detA / det
	qr.b = detB / det
	qr.c = detC / det
	qr.fitted = true

	// 結果を表示
	fmt.Println("=== 計算結果 ===")
	fmt.Printf("近似式: y = %.4fx² + %.4fx + %.4f\n", qr.a, qr.b, qr.c)
	fmt.Println()

	// 各点での予測値と誤差を計算
	fmt.Println("予測値と誤差:")
	var sumSquaredError float64
	for i, point := range qr.data {
		predicted := qr.Predict(point.X)
		error := point.Y - predicted
		sumSquaredError += error * error
		fmt.Printf("  点%d: x=%.2f, 実測値=%.2f, 予測値=%.4f, 誤差=%.4f\n",
			i+1, point.X, point.Y, predicted, error)
	}

	// 決定係数 R² を計算
	meanY := sumY / n
	var totalSumSquares float64
	for _, point := range qr.data {
		totalSumSquares += (point.Y - meanY) * (point.Y - meanY)
	}

	r2 := 1.0 - (sumSquaredError / totalSumSquares)

	fmt.Println()
	fmt.Printf("二乗誤差の和: %.6f\n", sumSquaredError)
	fmt.Printf("決定係数 R²: %.6f\n", r2)
	fmt.Println("(R²が1に近いほど、データへの当てはまりが良い)")
	fmt.Println()
}

// Predict は与えられたxに対してy値を予測する
func (qr *QuadraticRegression) Predict(x float64) float64 {
	if !qr.fitted {
		fmt.Println("警告: Fit()を先に実行してください")
		return 0
	}
	return qr.a*x*x + qr.b*x + qr.c
}

// VisualizeLeastSquares は最小二乗法の概念を視覚的に説明する
func VisualizeLeastSquares() {
	fmt.Println("=== 最小二乗法の原理 ===")
	fmt.Println("データ点と近似曲線の「誤差の二乗和」を最小化する方法")
	fmt.Println()
	fmt.Println("  y |")
	fmt.Println("    |    ●  <- データ点")
	fmt.Println("    |   /|")
	fmt.Println("    |  / | <- 誤差")
	fmt.Println("    | /  |")
	fmt.Println("    |/___●______ 近似直線")
	fmt.Println("    |    ")
	fmt.Println("    |  ●")
	fmt.Println("    |____________________ x")
	fmt.Println()
	fmt.Println("誤差 = 実測値 - 予測値")
	fmt.Println("目標: Σ(誤差²) を最小化")
	fmt.Println()
	fmt.Println("【1次式近似】")
	fmt.Println("  y = ax + b")
	fmt.Println("  正規方程式を解いて a, b を求める")
	fmt.Println()
	fmt.Println("【2次式近似】")
	fmt.Println("  y = ax² + bx + c")
	fmt.Println("  連立方程式（正規方程式）を解いて a, b, c を求める")
	fmt.Println()
}

func main() {
	fmt.Println("最小二乗法（Least Squares Method）のデモ")
	fmt.Println("==========================================\n")

	// 最小二乗法の原理を説明
	VisualizeLeastSquares()

	// 例1: 1次式近似（ほぼ線形のデータ）
	fmt.Println("\n【例1】1次式近似 - 線形に近いデータ")
	fmt.Println("真の関数: y = 2x + 3 (に近いデータ)")
	linearData := []DataPoint{
		{X: 1.0, Y: 5.1},
		{X: 2.0, Y: 7.0},
		{X: 3.0, Y: 8.9},
		{X: 4.0, Y: 11.1},
		{X: 5.0, Y: 12.8},
		{X: 6.0, Y: 15.2},
	}

	lr := NewLinearRegression(linearData)
	lr.Fit()

	// 新しいxで予測
	fmt.Println("新しい値の予測:")
	testX := []float64{2.5, 7.0}
	for _, x := range testX {
		predicted := lr.Predict(x)
		fmt.Printf("  x=%.1f のとき、y=%.4f と予測\n", x, predicted)
	}

	// 例2: 2次式近似（放物線のデータ）
	fmt.Println("\n【例2】2次式近似 - 放物線のデータ")
	fmt.Println("真の関数: y = x² - 2x + 1 (に近いデータ)")
	quadraticData := []DataPoint{
		{X: -2.0, Y: 9.1},
		{X: -1.0, Y: 4.2},
		{X: 0.0, Y: 0.9},
		{X: 1.0, Y: 0.1},
		{X: 2.0, Y: 1.0},
		{X: 3.0, Y: 3.8},
		{X: 4.0, Y: 9.2},
	}

	qr := NewQuadraticRegression(quadraticData)
	qr.Fit()

	// 新しいxで予測
	fmt.Println("新しい値の予測:")
	testX2 := []float64{0.5, 2.5}
	for _, x := range testX2 {
		predicted := qr.Predict(x)
		fmt.Printf("  x=%.1f のとき、y=%.4f と予測\n", x, predicted)
	}

	// 例3: 線形データに2次式近似を適用（オーバーフィッティングの比較）
	fmt.Println("\n【例3】比較: 同じデータに1次式と2次式を適用")
	fmt.Println("線形データに対して:")

	lr2 := NewLinearRegression(linearData)
	lr2.Fit()

	qr2 := NewQuadraticRegression(linearData)
	qr2.Fit()

	fmt.Println("→ 線形データの場合、1次式近似の方が適切")
	fmt.Println("  (2次式は不必要に複雑で、過学習の可能性)")
}
