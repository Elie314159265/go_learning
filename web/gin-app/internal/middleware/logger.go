// Package middleware はHTTPリクエストの前処理・後処理を提供します
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware はカスタムログ出力を行うミドルウェアです
// リクエストとレスポンスの詳細情報をログに記録します
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエスト開始時刻
		startTime := time.Now()

		// リクエスト情報の取得
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 次のミドルウェア/ハンドラーを実行
		c.Next()

		// レスポンス後の処理
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// クエリパラメータの追加
		if raw != "" {
			path = path + "?" + raw
		}

		// ログの色分け（ステータスコードに基づく）
		statusColor := getStatusColor(statusCode)

		// ログ出力
		fmt.Printf("[GIN] %v | %s %3d %s | %13v | %15s | %-7s %s %s\n",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusColor,
			statusCode,
			resetColor(),
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}

// RequestIDMiddleware はリクエストごとにユニークなIDを生成します
// 分散システムでのリクエスト追跡に便利です
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストIDの生成（簡易版）
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		// コンテキストに設定
		c.Set("request_id", requestID)

		// レスポンスヘッダーにも追加
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// ステータスコードに応じた色コードを返す
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return green() // 成功
	case statusCode >= 300 && statusCode < 400:
		return white() // リダイレクト
	case statusCode >= 400 && statusCode < 500:
		return yellow() // クライアントエラー
	default:
		return red() // サーバーエラー
	}
}

// ANSIカラーコード
func green() string  { return "\033[97;42m" }
func white() string  { return "\033[90;47m" }
func yellow() string { return "\033[90;43m" }
func red() string    { return "\033[97;41m" }
func resetColor() string { return "\033[0m" }
