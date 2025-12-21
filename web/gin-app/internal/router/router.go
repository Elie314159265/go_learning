// Package router はアプリケーションのルーティング設定を提供します
// 全てのエンドポイント、ミドルウェア、ハンドラーをここで設定します
package router

import (
	"net/http"
	"time"

	"go_learning/web/gin-app/internal/config"
	"go_learning/web/gin-app/internal/database"
	"go_learning/web/gin-app/internal/handlers"
	"go_learning/web/gin-app/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter はGinルーターを設定し、全てのルートを登録します
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Ginのモードを設定（debug, release, test）
	gin.SetMode(cfg.Server.Mode)

	// ルーターの作成
	r := gin.New()

	// グローバルミドルウェアの設定
	r.Use(gin.Recovery())                          // パニック時の自動復旧
	r.Use(middleware.LoggerMiddleware())           // カスタムロガー
	r.Use(middleware.RequestIDMiddleware())        // リクエストID生成
	r.Use(middleware.CORSMiddleware())             // CORS設定

	// レートリミッターの設定（1分間に100リクエストまで）
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute)
	r.Use(rateLimiter.RateLimitMiddleware())

	// ハンドラーの初期化
	userHandler := handlers.NewUserHandler(db, cfg)
	productHandler := handlers.NewProductHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		// データベース接続チェック
		if err := database.Ping(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"database":    "connected",
			"version":     cfg.App.Version,
			"environment": cfg.App.Environment,
		})
	})

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to " + cfg.App.Name,
			"version": cfg.App.Version,
			"docs":    "/api/v1/docs",
		})
	})

	// API v1 グループ
	v1 := r.Group("/api/v1")
	{
		// 認証エンドポイント（公開）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register) // ユーザー登録
			auth.POST("/login", userHandler.Login)       // ログイン
		}

		// ユーザーエンドポイント
		users := v1.Group("/users")
		{
			// 認証が必要なエンドポイント
			users.Use(middleware.AuthMiddleware(cfg))
			users.GET("/profile", userHandler.GetProfile)       // 自分のプロフィール取得
			users.PUT("/profile", userHandler.UpdateProfile)    // プロフィール更新

			// 管理者のみアクセス可能
			admin := users.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("", userHandler.ListUsers)          // 全ユーザー一覧
				admin.GET("/:id", userHandler.GetUser)        // 特定ユーザー取得
				admin.DELETE("/:id", userHandler.DeleteUser)  // ユーザー削除
			}
		}

		// 商品エンドポイント
		products := v1.Group("/products")
		{
			// 公開エンドポイント（認証不要）
			products.GET("", productHandler.ListProducts)              // 商品一覧
			products.GET("/:id", productHandler.GetProduct)            // 商品詳細
			products.GET("/categories", productHandler.GetCategories)  // カテゴリー一覧

			// 管理者のみアクセス可能
			admin := products.Group("")
			admin.Use(middleware.AuthMiddleware(cfg))
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("", productHandler.CreateProduct)           // 商品作成
				admin.PUT("/:id", productHandler.UpdateProduct)        // 商品更新
				admin.DELETE("/:id", productHandler.DeleteProduct)     // 商品削除
			}
		}

		// 注文エンドポイント（全て認証が必要）
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.POST("", orderHandler.CreateOrder)                  // 注文作成
			orders.GET("", orderHandler.ListOrders)                    // 注文一覧
			orders.GET("/:id", orderHandler.GetOrder)                  // 注文詳細
			orders.POST("/:id/cancel", orderHandler.CancelOrder)       // 注文キャンセル

			// 管理者のみアクセス可能
			admin := orders.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.PATCH("/:id/status", orderHandler.UpdateOrderStatus) // ステータス更新
			}
		}

		// APIドキュメントエンドポイント
		v1.GET("/docs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API Documentation",
				"version": "1.0.0",
				"endpoints": gin.H{
					"auth": gin.H{
						"POST /api/v1/auth/register": "ユーザー登録",
						"POST /api/v1/auth/login":    "ログイン",
					},
					"users": gin.H{
						"GET /api/v1/users/profile":    "プロフィール取得（認証必要）",
						"PUT /api/v1/users/profile":    "プロフィール更新（認証必要）",
						"GET /api/v1/users":            "全ユーザー一覧（管理者のみ）",
						"GET /api/v1/users/:id":        "ユーザー詳細（管理者のみ）",
						"DELETE /api/v1/users/:id":     "ユーザー削除（管理者のみ）",
					},
					"products": gin.H{
						"GET /api/v1/products":              "商品一覧",
						"GET /api/v1/products/:id":          "商品詳細",
						"GET /api/v1/products/categories":   "カテゴリー一覧",
						"POST /api/v1/products":             "商品作成（管理者のみ）",
						"PUT /api/v1/products/:id":          "商品更新（管理者のみ）",
						"DELETE /api/v1/products/:id":       "商品削除（管理者のみ）",
					},
					"orders": gin.H{
						"POST /api/v1/orders":               "注文作成（認証必要）",
						"GET /api/v1/orders":                "注文一覧（認証必要）",
						"GET /api/v1/orders/:id":            "注文詳細（認証必要）",
						"POST /api/v1/orders/:id/cancel":    "注文キャンセル（認証必要）",
						"PATCH /api/v1/orders/:id/status":   "ステータス更新（管理者のみ）",
					},
				},
			})
		})
	}

	// 404エラーハンドリング
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "エンドポイントが見つかりません",
			"path":  c.Request.URL.Path,
		})
	})

	return r
}
