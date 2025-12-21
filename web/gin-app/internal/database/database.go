// Package database はデータベース接続とマイグレーション機能を提供します
// GORM を使用してPostgreSQLやMySQLなどのデータベースとの接続を管理します
package database

import (
	"fmt"
	"log"

	"go_learning/web/gin-app/internal/config"
	"go_learning/web/gin-app/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB はデータベース接続のグローバル変数です（必要に応じて使用）
var DB *gorm.DB

// NewConnection は新しいデータベース接続を作成します
// 設定に基づいて接続プールの設定も行います
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// PostgreSQL用のDSN（Data Source Name）を生成
	dsn := cfg.GetDatabaseDSN()

	// GORMのログレベルを設定
	// Silent: ログなし、Error: エラーのみ、Warn: 警告以上、Info: 全て
	logLevel := logger.Info
	if cfg.SSLMode == "disable" {
		logLevel = logger.Warn // 開発環境ではログを少なめに
	}

	// データベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		// PrepareStmt: プリペアドステートメントを使用してパフォーマンス向上
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続プールの設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("データベースプール設定エラー: %w", err)
	}

	// 最大オープン接続数を設定（デフォルト: 無制限）
	// 多すぎるとデータベースサーバーに負荷がかかる
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 最大アイドル接続数を設定（デフォルト: 2）
	// 接続の再利用により新規接続のオーバーヘッドを削減
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// 接続の最大ライフタイムを設定
	// 古い接続を定期的にクローズして接続の健全性を保つ
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// グローバル変数に保存（オプション）
	DB = db

	log.Println("データベース接続に成功しました")
	return db, nil
}

// AutoMigrate はデータベースのテーブルを自動で作成・更新します
// 開発環境で便利ですが、本番環境では専用のマイグレーションツールの使用を推奨
func AutoMigrate(db *gorm.DB) error {
	log.Println("データベースマイグレーションを開始します...")

	// ここに全てのモデルを登録
	// GORMが自動的にテーブルを作成・更新します
	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)

	if err != nil {
		return fmt.Errorf("マイグレーションエラー: %w", err)
	}

	log.Println("マイグレーションが完了しました")
	return nil
}

// Close はデータベース接続をクローズします
// アプリケーション終了時に呼び出してリソースを解放します
func (db *gorm.DB) Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping はデータベース接続が有効かどうかをチェックします
// ヘルスチェックエンドポイントで使用できます
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
