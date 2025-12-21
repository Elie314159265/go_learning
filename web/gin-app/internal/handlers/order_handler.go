// Package handlers はHTTPリクエストを処理するハンドラー関数を提供します
package handlers

import (
	"net/http"
	"strconv"

	"go_learning/web/gin-app/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderHandler は注文関連のハンドラーをまとめる構造体です
type OrderHandler struct {
	db *gorm.DB
}

// NewOrderHandler は新しいOrderHandlerを作成します
func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// CreateOrder は新しい注文を作成します
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	// トランザクション開始
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 注文の作成
	order := models.Order{
		UserID:          userID.(uint),
		Status:          models.OrderStatusPending,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注文の作成に失敗しました",
		})
		return
	}

	// 注文明細の作成
	var totalAmount float64
	for _, item := range req.Items {
		// 商品情報の取得
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error": "商品が見つかりません: " + strconv.Itoa(int(item.ProductID)),
			})
			return
		}

		// 在庫チェック
		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "在庫が不足しています: " + product.Name,
			})
			return
		}

		// 注文明細の作成
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price, // 注文時の価格を記録
		}
		orderItem.Subtotal = float64(orderItem.Quantity) * orderItem.Price

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "注文明細の作成に失敗しました",
			})
			return
		}

		// 在庫を減らす
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "在庫の更新に失敗しました",
			})
			return
		}

		totalAmount += orderItem.Subtotal
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	// 合計金額の更新
	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注文の更新に失敗しました",
		})
		return
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注文の確定に失敗しました",
		})
		return
	}

	// 注文明細を含めて取得
	h.db.Preload("OrderItems.Product").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "注文を作成しました",
		"order":   order,
	})
}

// ListOrders はユーザーの注文リストを取得します
// GET /api/v1/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// ページネーション
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	query := h.db.Model(&models.Order{})

	// 管理者以外は自分の注文のみ表示
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// 総数を取得
	var total int64
	query.Count(&total)

	// 注文を取得（注文明細と商品情報も含む）
	var orders []models.Order
	if err := query.
		Preload("OrderItems.Product").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注文の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetOrder は特定の注文情報を取得します
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var order models.Order
	query := h.db.Preload("OrderItems.Product").Preload("User")

	// 管理者以外は自分の注文のみ表示
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "注文が見つかりません",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus は注文ステータスを更新します（管理者のみ）
// PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.OrderUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	var order models.Order
	if err := h.db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "注文が見つかりません",
		})
		return
	}

	// ステータスの更新
	order.Status = req.Status
	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ステータスの更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注文ステータスを更新しました",
		"order":   order,
	})
}

// CancelOrder は注文をキャンセルします
// POST /api/v1/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var order models.Order
	query := h.db.Preload("OrderItems")

	// 管理者以外は自分の注文のみキャンセル可能
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "注文が見つかりません",
		})
		return
	}

	// キャンセル可能なステータスチェック
	if order.Status == models.OrderStatusShipped || order.Status == models.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "発送済みまたは配達完了の注文はキャンセルできません",
		})
		return
	}

	// トランザクション開始
	tx := h.db.Begin()

	// 在庫を戻す
	for _, item := range order.OrderItems {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err == nil {
			product.Stock += item.Quantity
			tx.Save(&product)
		}
	}

	// ステータスをキャンセルに変更
	order.Status = models.OrderStatusCancelled
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注文のキャンセルに失敗しました",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "注文をキャンセルしました",
		"order":   order,
	})
}
