// Package models はデータベースのテーブル構造を定義します
package models

import (
	"time"

	"gorm.io/gorm"
)

// Order は注文情報を表すモデルです
type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 外部キー: ユーザーID
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"` // リレーション

	OrderNumber string        `gorm:"uniqueIndex;not null;size:50" json:"order_number"` // 注文番号
	Status      string        `gorm:"size:20;default:'pending'" json:"status"`          // 注文ステータス
	TotalAmount float64       `gorm:"type:decimal(10,2)" json:"total_amount"`           // 合計金額

	// 配送情報
	ShippingAddress string   `gorm:"type:text" json:"shipping_address"`                // 配送先住所
	BillingAddress  string   `gorm:"type:text" json:"billing_address"`                 // 請求先住所

	// リレーション: 1つの注文は複数の注文明細を持つ
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

// OrderItem は注文明細を表すモデルです
type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// 外部キー
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	Order     Order          `gorm:"foreignKey:OrderID" json:"-"`              // リレーション

	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"` // リレーション

	Quantity  int            `gorm:"not null" json:"quantity"`                 // 数量
	Price     float64        `gorm:"not null;type:decimal(10,2)" json:"price"` // 単価（注文時の価格）
	Subtotal  float64        `gorm:"type:decimal(10,2)" json:"subtotal"`       // 小計
}

// OrderStatus は注文ステータスの定数です
const (
	OrderStatusPending   = "pending"    // 保留中
	OrderStatusConfirmed = "confirmed"  // 確認済み
	OrderStatusShipped   = "shipped"    // 発送済み
	OrderStatusDelivered = "delivered"  // 配達完了
	OrderStatusCancelled = "cancelled"  // キャンセル
)

// OrderCreateRequest は注文作成時のリクエストボディです
type OrderCreateRequest struct {
	Items           []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
	ShippingAddress string            `json:"shipping_address" binding:"required,min=10"`
	BillingAddress  string            `json:"billing_address" binding:"required,min=10"`
}

// OrderItemRequest は注文明細のリクエストです
type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// OrderUpdateStatusRequest は注文ステータス更新のリクエストです
type OrderUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}

// BeforeCreate は注文作成前に実行されるGORMフックです
// 注文番号を自動生成します
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderNumber == "" {
		// 注文番号の生成（タイムスタンプベース）
		o.OrderNumber = "ORD" + time.Now().Format("20060102150405")
	}
	return nil
}

// BeforeSave は注文明細保存前に実行されるGORMフックです
// 小計を自動計算します
func (oi *OrderItem) BeforeSave(tx *gorm.DB) error {
	oi.Subtotal = float64(oi.Quantity) * oi.Price
	return nil
}

// CalculateTotalAmount は注文の合計金額を計算します
func (o *Order) CalculateTotalAmount(tx *gorm.DB) error {
	var total float64
	for _, item := range o.OrderItems {
		total += item.Subtotal
	}
	o.TotalAmount = total
	return tx.Model(o).Update("total_amount", total).Error
}
