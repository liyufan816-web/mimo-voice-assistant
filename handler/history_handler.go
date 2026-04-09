package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/model"
	"github.com/liyufan816-web/mimo-voice-assistant/service"
	"github.com/liyufan816-web/mimo-voice-assistant/utils"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	historyService *service.HistoryService
	authService    *service.AuthService
}

// NewHistoryHandler 创建历史记录处理器
func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{
		historyService: service.NewHistoryService(),
		authService:    service.NewAuthService(),
	}
}

// GetUserHistory 获取用户的历史记录（分页）
func (h *HistoryHandler) GetUserHistory(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSizeStr := c.Query("page_size")
	pageSize := 20
	if pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 && parsedPageSize <= 100 {
			pageSize = parsedPageSize
		}
	}

	records, total, err := h.historyService.GetUserRecords(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result := make([]gin.H, len(records))
	for i, record := range records {
		result[i] = gin.H{
			"id":         record.ID,
			"text":       record.Text,
			"style":      record.Style,
			"voice":      record.Voice,
			"speed":      record.Speed,
			"request_id": record.RequestID,
			"status":     record.Status,
			"success":    record.Status == "success",
			"message":    record.Message,
			"duration":   record.Duration,
			"file_size":  record.FileSize,
			"created_at": record.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	utils.SuccessResponse(c, gin.H{
		"records":     result,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize-1)) / int64(pageSize),
	})
}

// GetRecordDetail 获取单条记录详情
func (h *HistoryHandler) GetRecordDetail(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的记录 ID")
		return
	}

	record, err := h.historyService.GetRecordByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// 检查权限：只能查看自己的记录
	if record.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "无权访问该记录")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"record": gin.H{
			"id":         record.ID,
			"text":       record.Text,
			"style":      record.Style,
			"voice":      record.Voice,
			"speed":      record.Speed,
			"request_id": record.RequestID,
			"status":     record.Status,
			"success":    record.Status == "success",
			"message":    record.Message,
			"duration":   record.Duration,
			"file_size":  record.FileSize,
			"created_at": record.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteRecord 删除单条记录
func (h *HistoryHandler) DeleteRecord(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的记录 ID")
		return
	}

	err = h.historyService.DeleteRecord(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "删除成功",
	})
}

// ClearUserHistory 清空用户的历史记录
func (h *HistoryHandler) ClearUserHistory(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	err := h.historyService.ClearUserHistory(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "清空成功",
	})
}

// ExportUserHistory 导出用户历史记录为 CSV
func (h *HistoryHandler) ExportUserHistory(c *gin.Context) {
	userID := h.authService.GetUserIDFromContext(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未认证")
		return
	}

	var records []*model.TTSRecord
	if err := config.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment;filename=history.csv")
	c.String(http.StatusOK, getCSVContent(records))
}

// getCSVContent 生成 CSV 内容
func getCSVContent(records []*model.TTSRecord) string {
	const header = "ID,Text,Style,Voice,Speed,Status,CreatedAt\n"
	var sb strings.Builder
	sb.WriteString(header)

	for _, r := range records {
		text := strings.ReplaceAll(r.Text, "\"", "\"\"")
		sb.WriteString(fmt.Sprintf("%d,\"%s\",%s,%s,%.1f,%s,%s\n",
			r.ID, text, r.Style, r.Voice, r.Speed, r.Status, r.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	return sb.String()
}
