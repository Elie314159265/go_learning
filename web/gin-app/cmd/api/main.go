// Package main はアプリケーションのエントリーポイントです
// このファイルでは、設定の読み込み、データベース接続の初期化、
// ルーターのセットアップ、サーバーの起動を行います
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_learning/web/gin-app/internal/config"
	"go_learning/web/gin-app/internal/database"
	"go_learning/web/gin-app/internal/router"
)

func main() {
	// 1. 設定の読み込み
	// 環境変数や設定ファイルからアプリケーションの設定を読み込みます
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// 2. データベース接続の初期化
	// PostgreSQL等のデータベースに接続し、接続プールを作成します
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	defer database.Close(db)

	// 3. データベースマイグレーション
	// テーブルの作成や更新を自動で行います
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	// 4. ルーターのセットアップ
	// Ginのルーターを作成し、全てのエンドポイントとミドルウェアを設定します
	r := router.SetupRouter(db, cfg)

	// 5. HTTPサーバーの作成
	// タイムアウトやポート設定を含むHTTPサーバーを構成します
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 6. サーバーをゴルーチンで起動
	// メインゴルーチンをブロックせずにサーバーを起動します
	go func() {
		log.Printf("サーバーを起動します: http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバーの起動に失敗しました: %v", err)
		}
	}()

	// 7. グレースフルシャットダウンの設定
	// SIGINT (Ctrl+C) や SIGTERM シグナルを受信したときに、
	// 既存のリクエストを処理完了してから安全にサーバーを停止します
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("サーバーをシャットダウンしています...")

	// 8. シャットダウンのタイムアウト設定
	// 最大5秒間、既存のリクエストの完了を待ちます
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("サーバーの強制シャットダウン:", err)
	}

	log.Println("サーバーが正常に停止しました")
}
