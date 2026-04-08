package model

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateAPIKey 生成 API Key
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "mimo_" + base64.StdEncoding.EncodeToString(b), nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserIDKey 用户 ID 键
	UserIDKey ContextKey = "user_id"
	// UserKey 用户信息键
	UserKey ContextKey = "user"
	// APIKeyKey API 密钥键
	APIKeyKey ContextKey = "api_key"
)

// ExtractAPIKey 从请求中提取 API Key
func ExtractAPIKey(header string) string {
	// 支持 Bearer token 格式
	if strings.HasPrefix(header, "Bearer ") {
		return header[7:]
	}
	return header
}
