// Package models はデータベースのテーブル構造を定義します
package models

import (
	"time"

	"gorm.io/gorm"
)

// Product は商品情報を表すモデルです
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string         `gorm:"not null;size:200" json:"name"`            // 商品名
	Description string         `gorm:"type:text" json:"description"`             // 商品説明
	Price       float64        `gorm:"not null;type:decimal(10,2)" json:"price"` // 価格（小数点2桁まで）
	Stock       int            `gorm:"not null;default:0" json:"stock"`          // 在庫数
	SKU         string         `gorm:"uniqueIndex;size:50" json:"sku"`           // 商品コード（一意）
	Category    string         `gorm:"size:50" json:"category"`                  // カテゴリー
	ImageURL    string         `gorm:"size:500" json:"image_url"`                // 商品画像URL
	IsActive    bool           `gorm:"default:true" json:"is_active"`            // 販売中フラグ

	// リレーション: 商品は複数の注文明細に含まれる
	OrderItems  []OrderItem    `gorm:"foreignKey:ProductID" json:"-"`
}

// ProductCreateRequest は商品作成時のリクエストボディです
type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`        // 必須
	Description string  `json:"description" binding:"max=1000"`               // オプション
	Price       float64 `json:"price" binding:"required,gt=0"`                // 必須、0より大きい
	Stock       int     `json:"stock" binding:"required,gte=0"`               // 必須、0以上
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`          // 必須
	Category    string  `json:"category" binding:"max=50"`                    // オプション
	ImageURL    string  `json:"image_url" binding:"omitempty,url,max=500"`    // オプション、URL形式
}

// ProductUpdateRequest は商品更新時のリクエストボディです
type ProductUpdateRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=1,max=200"`
	Description string   `json:"description" binding:"max=1000"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`               // ポインタでnull許可
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	Category    string   `json:"category" binding:"max=50"`
	ImageURL    string   `json:"image_url" binding:"omitempty,url,max=500"`
	IsActive    *bool    `json:"is_active"`
}

// BeforeSave は保存前に実行されるGORMフックです
// 在庫が0の場合は自動的に非アクティブにします
func (p *Product) BeforeSave(tx *gorm.DB) error {
	if p.Stock == 0 {
		p.IsActive = false
	}
	return nil
}
