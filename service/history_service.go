package service

import (
	"fmt"
	"time"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/model"
)

// HistoryService 历史记录服务
type HistoryService struct{}

// NewHistoryService 创建历史记录服务
func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

// CreateRecord 创建 TTS 记录
func (s *HistoryService) CreateRecord(userID uint, text, style, voice string, requestID string, success bool, errorMsg string, duration int, fileSize int64) (*model.TTSRecord, error) {
	record := &model.TTSRecord{
		UserID:    userID,
		Text:      text,
		Style:     style,
		Voice:     voice,
		RequestID: requestID,
		Status:    "pending", // 使用 Status 而不是 Success
		Message:   errorMsg,
		Duration:  duration,
		FileSize:  fileSize,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if success {
		record.Status = "success"
	} else {
		record.Status = "error"
	}

	if err := config.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("创建历史记录失败：%v", err)
	}

	return record, nil
}

// GetUserRecords 获取用户的 TTS 记录（分页）
func (s *HistoryService) GetUserRecords(userID uint, page, pageSize int) ([]*model.TTSRecord, int64, error) {
	var records []*model.TTSRecord
	var total int64

	offset := (page - 1) * pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	if err := config.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询历史记录失败：%v", err)
	}

	config.DB.Model(&model.TTSRecord{}).Where("user_id = ?", userID).Count(&total)

	return records, total, nil
}

// GetRecordByID 根据 ID 获取单条记录
func (s *HistoryService) GetRecordByID(id uint) (*model.TTSRecord, error) {
	var record model.TTSRecord
	if err := config.DB.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("记录不存在")
	}
	return &record, nil
}

// DeleteRecord 删除记录
func (s *HistoryService) DeleteRecord(id uint) error {
	if err := config.DB.Delete(&model.TTSRecord{}, id).Error; err != nil {
		return fmt.Errorf("删除记录失败：%v", err)
	}
	return nil
}

// ClearUserHistory 清空用户的历史记录
func (s *HistoryService) ClearUserHistory(userID uint) error {
	if err := config.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Delete(&model.TTSRecord{}).Error; err != nil {
		return fmt.Errorf("清空历史记录失败：%v", err)
	}
	return nil
}

// GetRecentRecords 获取最近的 N 条记录（不限制用户，用于统计等）
func (s *HistoryService) GetRecentRecords(limit int) ([]*model.TTSRecord, error) {
	var records []*model.TTSRecord
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	if err := config.DB.Order("created_at DESC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("查询最近记录失败：%v", err)
	}

	return records, nil
}

// RecordSuccess 记录成功的 TTS 请求
func (s *HistoryService) RecordSuccess(userID, requestID uint, text, style, voice string) error {
	now := time.Now()
	// 尝试查找已存在的记录
	var record model.TTSRecord
	if err := config.DB.Where("request_id = ?", requestID).First(&record).Error; err == nil {
		record.Status = "success" // 修改此处，使用 Status 而不是 Success
		record.UpdatedAt = now
		return config.DB.Save(&record).Error
	}

	// 创建新记录
	newRecord := &model.TTSRecord{
		UserID:    userID,
		Text:      text,
		Style:     style,
		Voice:     voice,
		RequestID: fmt.Sprintf("%d", requestID),
		Status:    "success",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return config.DB.Create(newRecord).Error
}

// RecordError 记录失败的 TTS 请求
func (s *HistoryService) RecordError(requestID string, text, style, errorMsg string) error {
	now := time.Now()

	record := &model.TTSRecord{
		RequestID: requestID,
		Text:      text,
		Style:     style,
		Status:    "error",
		Message:   errorMsg,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return config.DB.Create(record).Error
}