// Package utils は汎用的なユーティリティ関数を提供します
package utils

import (
	"regexp"
	"unicode"
)

// IsValidEmail はメールアドレスの形式が正しいかを検証します
func IsValidEmail(email string) bool {
	// RFC 5322に準拠した簡易的なメールアドレス検証
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsStrongPassword はパスワードの強度をチェックします
// 少なくとも8文字以上で、大文字、小文字、数字を含む必要があります
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

// IsValidUsername はユーザー名の形式が正しいかを検証します
// 3〜20文字の英数字とアンダースコアのみ許可
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeInput は入力文字列から危険な文字を除去します
// XSS攻撃を防ぐための簡易的なサニタイズ
func SanitizeInput(input string) string {
	// HTMLタグを除去
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
