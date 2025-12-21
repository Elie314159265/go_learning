// Package middleware はHTTPリクエストの前処理・後処理を提供します
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter はレートリミット（アクセス制限）を実装する構造体です
// 特定IPアドレスからのリクエスト数を制限し、DDoS攻撃を防ぎます
type RateLimiter struct {
	visitors map[string]*Visitor // IPアドレスごとの訪問者情報
	mu       sync.RWMutex        // 並行アクセス制御用のミューテックス
	rate     int                 // 許可するリクエスト数
	duration time.Duration       // レート制限の期間
}

// Visitor は訪問者の情報を保持します
type Visitor struct {
	lastSeen time.Time // 最終アクセス時刻
	count    int       // リクエスト数
}

// NewRateLimiter は新しいRateLimiterを作成します
// rate: 許可するリクエスト数
// duration: レート制限の期間（例: 1分間に10リクエストなら rate=10, duration=1*time.Minute）
func NewRateLimiter(rate int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		duration: duration,
	}

	// 定期的に古い訪問者情報をクリーンアップ
	go rl.cleanupVisitors()

	return rl
}

// RateLimitMiddleware はレート制限を適用するミドルウェアです
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// レート制限チェック
		if !rl.allowRequest(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "リクエスト数が多すぎます。しばらく待ってから再試行してください。",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// allowRequest はリクエストを許可するかどうかを判定します
func (rl *RateLimiter) allowRequest(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 訪問者情報の取得または作成
	visitor, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &Visitor{
			lastSeen: now,
			count:    1,
		}
		return true
	}

	// レート制限期間内かチェック
	if now.Sub(visitor.lastSeen) > rl.duration {
		// 期間が過ぎていればカウントをリセット
		visitor.count = 1
		visitor.lastSeen = now
		return true
	}

	// レート制限チェック
	if visitor.count >= rl.rate {
		return false
	}

	// カウントを増やす
	visitor.count++
	return true
}

// cleanupVisitors は一定時間アクセスがない訪問者情報を削除します
// メモリリークを防ぐため、バックグラウンドで定期実行されます
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for ip, visitor := range rl.visitors {
			// 10分以上アクセスがない場合は削除
			if now.Sub(visitor.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}

		rl.mu.Unlock()
	}
}
