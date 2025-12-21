// Package handlers はHTTPリクエストを処理するハンドラー関数を提供します
package handlers

import (
	"net/http"
	"strconv"

	"go_learning/web/gin-app/internal/config"
	"go_learning/web/gin-app/internal/models"
	"go_learning/web/gin-app/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler はユーザー関連のハンドラーをまとめる構造体です
type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewUserHandler は新しいUserHandlerを作成します
func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:  db,
		cfg: cfg,
	}
}

// Register はユーザー登録を処理します
// POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest

	// リクエストボディのバリデーション
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	// ユーザー名の重複チェック
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "このユーザー名は既に使用されています",
		})
		return
	}

	// メールアドレスの重複チェック
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "このメールアドレスは既に使用されています",
		})
		return
	}

	// ユーザーの作成
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // BeforeCreateフックで自動的にハッシュ化されます
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user", // デフォルトはuserロール
		IsActive:  true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーの作成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ユーザー登録が完了しました",
		"user":    user.ToResponse(),
	})
}

// Login はユーザーログインを処理します
// POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	// ユーザーの検索（ユーザー名またはメールアドレス）
	var user models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ユーザー名またはパスワードが正しくありません",
		})
		return
	}

	// パスワードの検証
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ユーザー名またはパスワードが正しくありません",
		})
		return
	}

	// アクティブユーザーのチェック
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "このアカウントは無効化されています",
		})
		return
	}

	// JWTトークンの生成
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, h.cfg.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トークンの生成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ログインに成功しました",
		"token":   token,
		"user":    user.ToResponse(),
	})
}

// GetProfile は認証済みユーザーのプロフィールを取得します
// GET /api/v1/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// ミドルウェアで設定されたユーザーIDを取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ユーザーが見つかりません",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateProfile はユーザープロフィールを更新します
// PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力値が無効です: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ユーザーが見つかりません",
		})
		return
	}

	// 更新するフィールドのみ適用
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "プロフィールの更新に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "プロフィールを更新しました",
		"user":    user.ToResponse(),
	})
}

// ListUsers は全ユーザーのリストを取得します（管理者のみ）
// GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// ページネーション
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	// 総数を取得
	h.db.Model(&models.User{}).Count(&total)

	// ページネーション付きでユーザーを取得
	if err := h.db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーの取得に失敗しました",
		})
		return
	}

	// レスポンスに変換
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      userResponses,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetUser は特定のユーザー情報を取得します（管理者のみ）
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ユーザーが見つかりません",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser はユーザーを削除します（ソフトデリート）
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーの削除に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザーを削除しました",
	})
}
