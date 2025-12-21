// Package handlers はHTTPリクエストを処理するハンドラー関数を提供します
package handlers

import (
	"net/http"
	"strconv"

	"go_learning/web/gin-app/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler は商品関連のハンドラーをまとめる構造体です
type ProductHandler struct {
	db *gorm.DB
}

// NewProductHandler は新しいProductHandlerを作成します
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// CreateProduct は新しい商品を作成します（管理者のみ）
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	// SKUの重複チェック
	var existingProduct models.Product
	if err := h.db.Where("sku = ?", req.SKU).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "このSKUは既に使用されています",
		})
		return
	}

	// 商品の作成
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "商品の作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "商品を作成しました",
		"product": product,
	})
}

// ListProducts は商品リストを取得します（公開API）
// GET /api/v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// クエリパラメータ
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	searchQuery := c.Query("search")
	activeOnly := c.DefaultQuery("active_only", "true") == "true"

	offset := (page - 1) * pageSize

	// クエリの構築
	query := h.db.Model(&models.Product{})

	// アクティブな商品のみ表示（管理者以外）
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// カテゴリーフィルター
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 検索フィルター（商品名と説明で検索）
	if searchQuery != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// 総数を取得
	var total int64
	query.Count(&total)

	// 商品を取得
	var products []models.Product
	if err := query.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "商品の取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":    products,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetProduct は特定の商品情報を取得します
// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "商品が見つかりません",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct は商品情報を更新します（管理者のみ）
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "商品が見つかりません",
		})
		return
	}

	// 更新するフィールドのみ適用
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "商品の更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "商品を更新しました",
		"product": product,
	})
}

// DeleteProduct は商品を削除します（ソフトデリート、管理者のみ）
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "商品の削除に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "商品を削除しました",
	})
}

// GetCategories は商品カテゴリーのリストを取得します
// GET /api/v1/products/categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	var categories []string

	// DISTINCT でカテゴリーを取得
	if err := h.db.Model(&models.Product{}).
		Distinct("category").
		Where("category != ''").
		Pluck("category", &categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "カテゴリーの取得に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
