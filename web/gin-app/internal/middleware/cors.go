// Package middleware はHTTPリクエストの前処理・後処理を提供します
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CORSMiddleware はCORS（Cross-Origin Resource Sharing）の設定を行います
// フロントエンドアプリケーションが異なるドメインからAPIを呼び出せるようにします
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 許可するオリジン（開発環境の例）
		// 本番環境では具体的なドメインを指定することを推奨
		AllowOrigins: []string{
			"http://localhost:3000",  // React等のフロントエンド開発サーバー
			"http://localhost:8080",  // 同一ポート
			"http://localhost:5173",  // Vite等
		},

		// 許可するHTTPメソッド
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		// 許可するリクエストヘッダー
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},

		// レスポンスで公開するヘッダー
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},

		// クレデンシャル（Cookie等）の送信を許可
		AllowCredentials: true,

		// プリフライトリクエストのキャッシュ時間
		MaxAge: 12 * time.Hour,
	})
}

// CustomCORSMiddleware はカスタムCORS設定を提供します
// より細かい制御が必要な場合に使用します
func CustomCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 許可するオリジンのリスト
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:8080": true,
			// 本番環境のドメインを追加
			// "https://yourdomain.com": true,
		}

		// オリジンのチェック
		if _, ok := allowedOrigins[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		// プリフライトリクエスト（OPTIONS）の処理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
