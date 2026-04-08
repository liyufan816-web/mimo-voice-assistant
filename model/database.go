package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 不返回密码
	APIKey    string         `gorm:"uniqueIndex;size:64" json:"api_key"` // API 访问密钥
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// HistoryItem 历史记录模型
type HistoryItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Text      string         `gorm:"type:text;not null" json:"text"`
	Style     string         `gorm:"size:100" json:"style"`
	Speed     float64        `gorm:"default:1.0" json:"speed"`
	RequestID string         `gorm:"size:64" json:"request_id"`
	Success   bool           `gorm:"default:true" json:"success"`
	ErrorMsg  string         `gorm:"type:text" json:"error_msg"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HistoryItem) TableName() string {
	return "history_items"
}

// TTSRecord TTS 记录模型（完整记录）
type TTSRecord struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Text      string         `gorm:"type:text;not null" json:"text"`
	Style     string         `gorm:"size:100" json:"style"`
	Voice     string         `gorm:"size:50" json:"voice"`
	Speed     float64        `gorm:"default:1.0" json:"speed"`
	RequestID string         `gorm:"size:64;index" json:"request_id"`
	Status    string         `gorm:"size:20;default:'pending'" json:"status"` // pending, success, error
	Message   string         `gorm:"type:text" json:"message"`
	Duration  int            `gorm:"default:0" json:"duration"` // 音频时长（秒）
	FileSize  int64          `gorm:"default:0" json:"file_size"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TTSRecord) TableName() string {
	return "tts_records"
}

// PageRequest 分页请求参数
type PageRequest struct {
	Page  int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (r *PageRequest) GetOffset() int {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 20
	}
	return (r.Page - 1) * r.PageSize
}