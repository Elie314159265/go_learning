// Package utils は汎用的なユーティリティ関数を提供します
package utils

import (
	"errors"
	"time"

	"go_learning/web/gin-app/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims はJWTトークンに含まれる情報（クレーム）を定義します
type JWTClaims struct {
	UserID   uint   `json:"user_id"`   // ユーザーID
	Username string `json:"username"`  // ユーザー名
	Role     string `json:"role"`      // ロール（user, admin等）
	jwt.RegisteredClaims                // 標準クレーム（exp, iat等）
}

// GenerateJWT はJWTトークンを生成します
// ユーザーの認証情報を含む署名付きトークンを返します
func GenerateJWT(userID uint, username, role string, cfg config.JWTConfig) (string, error) {
	// クレームの作成
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expiration)), // 有効期限
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 発行時刻
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 有効開始時刻
			Issuer:    cfg.Issuer,                                         // 発行者
			Subject:   username,                                           // サブジェクト（ユーザー名）
		},
	}

	// トークンの作成（HS256アルゴリズムを使用）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名してトークン文字列を生成
	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT はJWTトークンを検証し、クレームを返します
// トークンが無効または期限切れの場合はエラーを返します
func ValidateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	// トークンのパース
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無効な署名アルゴリズムです")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// クレームの取得
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("無効なトークンです")
}

// RefreshToken は既存のトークンから新しいトークンを生成します
// 有効期限が近いトークンをリフレッシュする際に使用します
func RefreshToken(oldToken string, cfg config.JWTConfig) (string, error) {
	// 古いトークンの検証（有効期限切れでもクレームは取得）
	claims, err := ValidateJWT(oldToken, cfg.SecretKey)
	if err != nil {
		// 期限切れの場合も許可（リフレッシュのため）
		token, _ := jwt.ParseWithClaims(oldToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.SecretKey), nil
		})

		if token != nil {
			if c, ok := token.Claims.(*JWTClaims); ok {
				claims = c
			} else {
				return "", errors.New("トークンのクレームが取得できません")
			}
		} else {
			return "", errors.New("無効なトークンです")
		}
	}

	// 新しいトークンを生成
	return GenerateJWT(claims.UserID, claims.Username, claims.Role, cfg)
}
