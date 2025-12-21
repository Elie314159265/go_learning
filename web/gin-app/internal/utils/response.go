// Package utils は汎用的なユーティリティ関数を提供します
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse は標準的なエラーレスポンスの構造です
type ErrorResponse struct {
	Error   string      `json:"error"`             // エラーメッセージ
	Details interface{} `json:"details,omitempty"` // 詳細情報（オプション）
	Code    string      `json:"code,omitempty"`    // エラーコード（オプション）
}

// SuccessResponse は標準的な成功レスポンスの構造です
type SuccessResponse struct {
	Message string      `json:"message"`          // 成功メッセージ
	Data    interface{} `json:"data,omitempty"`   // データ（オプション）
}

// RespondWithError はエラーレスポンスを返します
func RespondWithError(c *gin.Context, statusCode int, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Error:   message,
		Details: details,
	})
}

// RespondWithSuccess は成功レスポンスを返します
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// BadRequest は400 Bad Requestレスポンスを返します
func BadRequest(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message, nil)
}

// Unauthorized は401 Unauthorizedレスポンスを返します
func Unauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message, nil)
}

// Forbidden は403 Forbiddenレスポンスを返します
func Forbidden(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, message, nil)
}

// NotFound は404 Not Foundレスポンスを返します
func NotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message, nil)
}

// InternalServerError は500 Internal Server Errorレスポンスを返します
func InternalServerError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message, nil)
}

// Created は201 Createdレスポンスを返します
func Created(c *gin.Context, message string, data interface{}) {
	RespondWithSuccess(c, http.StatusCreated, message, data)
}

// OK は200 OKレスポンスを返します
func OK(c *gin.Context, message string, data interface{}) {
	RespondWithSuccess(c, http.StatusOK, message, data)
}
