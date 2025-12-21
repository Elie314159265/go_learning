// Package config はアプリケーションの設定管理を提供します
// 環境変数や設定ファイルから設定を読み込み、アプリケーション全体で使用できるようにします
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション全体の設定を保持する構造体です
type Config struct {
	Server   ServerConfig   // サーバー関連の設定
	Database DatabaseConfig // データベース関連の設定
	JWT      JWTConfig      // JWT認証の設定
	App      AppConfig      // アプリケーション全般の設定
}

// ServerConfig はHTTPサーバーの設定を保持します
type ServerConfig struct {
	Port         string        // サーバーがリッスンするポート番号
	Mode         string        // Ginの動作モード (debug, release, test)
	ReadTimeout  time.Duration // リクエスト読み込みのタイムアウト
	WriteTimeout time.Duration // レスポンス書き込みのタイムアウト
	IdleTimeout  time.Duration // アイドル接続のタイムアウト
}

// DatabaseConfig はデータベース接続の設定を保持します
type DatabaseConfig struct {
	Host            string // データベースのホスト名
	Port            string // データベースのポート番号
	User            string // データベースのユーザー名
	Password        string // データベースのパスワード
	DBName          string // データベース名
	SSLMode         string // SSL接続モード (disable, require, verify-ca, verify-full)
	MaxOpenConns    int    // 最大オープン接続数
	MaxIdleConns    int    // 最大アイドル接続数
	ConnMaxLifetime time.Duration // 接続の最大ライフタイム
}

// JWTConfig はJWT認証の設定を保持します
type JWTConfig struct {
	SecretKey       string        // JWT署名用の秘密鍵
	ExpirationHours int           // トークンの有効期限（時間）
	Issuer          string        // トークンの発行者
	Expiration      time.Duration // トークンの有効期限（Duration）
}

// AppConfig はアプリケーション全般の設定を保持します
type AppConfig struct {
	Name        string // アプリケーション名
	Version     string // アプリケーションのバージョン
	Environment string // 実行環境 (development, staging, production)
	LogLevel    string // ログレベル (debug, info, warn, error)
}

// Load は環境変数から設定を読み込み、Config構造体を返します
// 環境変数が設定されていない場合は、デフォルト値を使用します
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "gin_app"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
			ExpirationHours: getIntEnv("JWT_EXPIRATION_HOURS", 24),
			Issuer:          getEnv("JWT_ISSUER", "gin-app"),
			Expiration:      time.Duration(getIntEnv("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Gin Web Application"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "debug"),
		},
	}

	// 必須の環境変数のバリデーション
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate は設定の妥当性をチェックします
func (c *Config) Validate() error {
	// JWT秘密鍵がデフォルトのままでproduction環境の場合はエラー
	if c.App.Environment == "production" &&
	   c.JWT.SecretKey == "your-secret-key-change-this-in-production" {
		return fmt.Errorf("production環境ではJWT_SECRETを必ず変更してください")
	}

	// データベース接続情報の基本チェック
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOSTが設定されていません")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAMEが設定されていません")
	}

	return nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv は環境変数を整数として取得し、存在しないまたは変換できない場合はデフォルト値を返します
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getDurationEnv は環境変数をDurationとして取得し、存在しないまたは変換できない場合はデフォルト値を返します
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetDatabaseDSN はPostgreSQL接続文字列を生成します
func (c *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
