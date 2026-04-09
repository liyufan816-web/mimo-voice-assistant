package service

import (
	"fmt"
	"strings"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/model"

	"github.com/gin-gonic/gin"
)

// AuthResult 认证结果
type AuthResult struct {
	User      *model.User
	Token     string
	Error     error
	IsAPIAuth bool
}

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 用户注册
func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	if email != "" {
		var existingEmail model.User
		if err := config.DB.Where("email = ?", email).First(&existingEmail).Error; err == nil {
			return nil, fmt.Errorf("邮箱已被使用")
		}
	}

	// 加密密码
	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败：%v", err)
	}

	// 生成 API Key
	apiKey, err := model.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("生成 API Key 失败：%v", err)
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		APIKey:   apiKey,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败：%v", err)
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*model.User, error) {
	var user model.User
	if err := config.DB.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	if !model.CheckPassword(password, user.Password) {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	return &user, nil
}

// GetUserIDFromContext 从上下文中获取用户 ID
func (s *AuthService) GetUserIDFromContext(c *gin.Context) uint {
	userID, ok := c.Get(string(model.UserIDKey)) // 使用 string() 而不是 String()
	if !ok {
		return 0
	}

	switch v := userID.(type) {
	case uint:
		return v
	case uint64:
		return uint(v)
	case int:
		return uint(v)
	case int64:
		return uint(v)
	case float64:
		return uint(v)
	default:
		return 0
	}
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return &user, nil
}

// APIKey 格式常量
const apiKeyPrefix = "mimo_"

// IsAPIKey 判断是否是有效的 API Key 格式
func IsAPIKey(key string) bool {
	return strings.HasPrefix(key, apiKeyPrefix) && len(key) > len(apiKeyPrefix)
}

// ValidateAPIKey 验证 API Key
func (s *AuthService) ValidateAPIKey(apiKey string) (*model.User, error) {
	if !IsAPIKey(apiKey) {
		return nil, fmt.Errorf("无效的 API Key 格式")
	}

	var user model.User
	if err := config.DB.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("API Key 无效")
	}

	return &user, nil
}