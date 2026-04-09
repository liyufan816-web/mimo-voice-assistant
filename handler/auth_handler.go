package handler

import (
	"net/http"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/model"
	"github.com/liyufan816-web/mimo-voice-assistant/service"
	"github.com/liyufan816-web/mimo-voice-assistant/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效："+err.Error())
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"api_key":    user.APIKey,
			"created_at": user.CreatedAt,
		},
		"message": "注册成功",
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效："+err.Error())
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"api_key":    user.APIKey,
			"created_at": user.CreatedAt,
		},
		"message": "登录成功",
	})
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

// GetAPIKey 获取用户的 API Key（需要登录后）
func (h *AuthHandler) GetAPIKey(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"api_key": user.APIKey,
	})
}

// Middleware 认证中间件
func (h *AuthHandler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 先检查 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			apiKey := model.ExtractAPIKey(authHeader)
			if apiKey != "" {
				user, err := h.authService.ValidateAPIKey(apiKey)
				if err == nil {
					c.Set(string(model.UserIDKey), user.ID)
					c.Set(string(model.UserKey), user)
					c.Set(string(model.APIKeyKey), apiKey)
					c.Next()
					return
				}
			}
		}

		// 2. 检查 Form 字段或查询参数中的 api_key
		apiKeyForm := c.PostForm("api_key")
		if apiKeyForm == "" {
			apiKeyForm = c.Query("api_key")
		}
		if apiKeyForm != "" {
			user, err := h.authService.ValidateAPIKey(apiKeyForm)
			if err == nil {
				c.Set(string(model.UserIDKey), user.ID)
				c.Set(string(model.UserKey), user)
				c.Set(string(model.APIKeyKey), apiKeyForm)
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusUnauthorized, "未提供有效的认证凭据")
		c.Abort()
	}
}

// CheckUserExists 检查用户名是否已存在
func (h *AuthHandler) CheckUserExists(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		username = c.Query("email")
	}

	if username == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "请提供 username 或 email")
		return
	}

	var count int64
	config.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, username).Count(&count)

	if count > 0 {
		utils.SuccessResponse(c, gin.H{
			"exists": true,
		})
	} else {
		utils.SuccessResponse(c, gin.H{
			"exists": false,
		})
	}
}