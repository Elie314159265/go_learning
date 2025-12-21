// Package models はデータベースのテーブル構造を定義します
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User はユーザー情報を表すモデルです
// GORMのタグを使用してデータベースのカラム設定を定義します
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                       // 主キー（自動インクリメント）
	CreatedAt time.Time      `json:"created_at"`                                 // 作成日時（GORM自動管理）
	UpdatedAt time.Time      `json:"updated_at"`                                 // 更新日時（GORM自動管理）
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`         // 削除日時（ソフトデリート用）

	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"` // ユーザー名（一意制約）
	Email     string         `gorm:"uniqueIndex;not null;size:100" json:"email"`   // メールアドレス（一意制約）
	Password  string         `gorm:"not null;size:255" json:"-"`                   // パスワード（JSON出力から除外）
	FirstName string         `gorm:"size:50" json:"first_name"`                    // 名
	LastName  string         `gorm:"size:50" json:"last_name"`                     // 姓
	Role      string         `gorm:"size:20;default:'user'" json:"role"`           // ロール（user, admin等）
	IsActive  bool           `gorm:"default:true" json:"is_active"`                // アクティブフラグ

	// リレーション: 1ユーザーは複数の注文を持つ
	Orders    []Order        `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// UserCreateRequest はユーザー作成時のリクエストボディです
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`        // 必須、3〜50文字
	Email     string `json:"email" binding:"required,email,max=100"`          // 必須、メール形式
	Password  string `json:"password" binding:"required,min=8,max=100"`       // 必須、8文字以上
	FirstName string `json:"first_name" binding:"max=50"`                     // オプション
	LastName  string `json:"last_name" binding:"max=50"`                      // オプション
}

// UserUpdateRequest はユーザー更新時のリクエストボディです
type UserUpdateRequest struct {
	Email     string `json:"email" binding:"omitempty,email,max=100"`         // オプション
	FirstName string `json:"first_name" binding:"max=50"`                     // オプション
	LastName  string `json:"last_name" binding:"max=50"`                      // オプション
	IsActive  *bool  `json:"is_active"`                                       // オプション（ポインタでnull許可）
}

// UserLoginRequest はログイン時のリクエストボディです
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`                     // ユーザー名またはメール
	Password string `json:"password" binding:"required"`                     // パスワード
}

// UserResponse はユーザー情報のレスポンスです（パスワードを除外）
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate はユーザー作成前に自動実行されるGORMフックです
// パスワードをハッシュ化します
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// パスワードが既にハッシュ化されている場合はスキップ
	if len(u.Password) > 0 && u.Password[0] != '$' {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword は入力されたパスワードが正しいかを検証します
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse はUserモデルをUserResponseに変換します
// パスワードなどの機密情報を除外してクライアントに返します
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
