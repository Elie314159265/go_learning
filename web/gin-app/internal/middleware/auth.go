// Package middleware はHTTPリクエストの前処理・後処理を提供します
package middleware

import (
	"net/http"
	"strings"

	"go_learning/web/gin-app/internal/config"
	"go_learning/web/gin-app/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware はJWTトークンによる認証を行うミドルウェアです
// リクエストヘッダーの Authorization: Bearer <token> からトークンを取得し、
// 検証に成功した場合はユーザー情報をコンテキストに設定します
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Authorizationヘッダーの取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "認証ヘッダーがありません",
			})
			c.Abort() // 以降の処理を中断
			return
		}

		// 2. "Bearer <token>" 形式のチェック
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無効な認証ヘッダー形式です",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. JWTトークンの検証
		claims, err := utils.ValidateJWT(tokenString, cfg.JWT.SecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "無効なトークンです: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 4. ユーザー情報をコンテキストに設定
		// ハンドラーでc.Get("user_id")等で取得可能になります
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 5. 次のミドルウェアまたはハンドラーを実行
		c.Next()
	}
}

// AdminMiddleware は管理者権限をチェックするミドルウェアです
// AuthMiddleware の後に実行する必要があります
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// コンテキストからロールを取得
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "認証が必要です",
			})
			c.Abort()
			return
		}

		// 管理者権限のチェック
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "管理者権限が必要です",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware はオプショナルな認証ミドルウェアです
// トークンがあれば検証し、なければ次の処理に進みます
// 公開/非公開コンテンツを同じエンドポイントで扱う場合に便利です
func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// トークンがない場合はそのまま次へ
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := utils.ValidateJWT(tokenString, cfg.JWT.SecretKey)
			if err == nil {
				// 有効なトークンの場合のみコンテキストに設定
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
			}
		}

		c.Next()
	}
}
