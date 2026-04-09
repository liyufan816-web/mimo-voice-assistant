package handler

import (
	"log"
	"net/http"

	"github.com/liyufan816-web/mimo-voice-assistant/model"
	"github.com/liyufan816-web/mimo-voice-assistant/service"
	"github.com/liyufan816-web/mimo-voice-assistant/utils"

	"github.com/gin-gonic/gin"
)

type TTSHandler struct {
	service *service.TTSService
}

func NewTTSHandler(service *service.TTSService) *TTSHandler {
	return &TTSHandler{service: service}
}

// 文本转语音
func (h *TTSHandler) ConvertText(c *gin.Context) {
	var req model.TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[REQUEST] POST /api/v1/tts/convert | 客户端IP: %s | 请求体解析失败: %v", c.ClientIP(), err)
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	log.Printf("[REQUEST] POST /api/v1/tts/convert | 客户端IP: %s | text长度: %d | style: %q", c.ClientIP(), len(req.Text), req.Style)

	if len(req.Text) == 0 {
		log.Printf("[RESPONSE] POST /api/v1/tts/convert | 状态: 400 | 文本不能为空")
		utils.ErrorResponse(c, http.StatusBadRequest, "文本不能为空")
		return
	}

	if len(req.Text) > 1000 {
		log.Printf("[RESPONSE] POST /api/v1/tts/convert | 状态: 400 | 文本长度超过1000")
		utils.ErrorResponse(c, http.StatusBadRequest, "文本长度不能超过1000字符")
		return
	}

	result, err := h.service.ConvertTextToSpeech(&req)
	if err != nil {
		log.Printf("[RESPONSE] POST /api/v1/tts/convert | 状态: 500 | 错误: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "语音转换失败: "+err.Error())
		return
	}

	if !result.Success {
		log.Printf("[RESPONSE] POST /api/v1/tts/convert | 状态: 400 | API返回: %s", result.Message)
		utils.ErrorResponse(c, http.StatusBadRequest, result.Message)
		return
	}

	log.Printf("[RESPONSE] POST /api/v1/tts/convert | 状态: 200 | 成功 | request_id: %s | 音频大小: %d bytes", result.RequestID, len(result.AudioData))
	c.Data(http.StatusOK, "audio/mpeg", result.AudioData)
}

// 健康检查
func (h *TTSHandler) HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"status":  "healthy",
		"service": "tts-backend",
	})
}

// 支持的音色列表
func (h *TTSHandler) GetVoices(c *gin.Context) {
	voices := []gin.H{
		{"id": "alloy", "name": "Alloy", "language": "zh-CN"},
		{"id": "echo", "name": "Echo", "language": "zh-CN"},
		{"id": "fable", "name": "Fable", "language": "zh-CN"},
		{"id": "onyx", "name": "Onyx", "language": "zh-CN"},
		{"id": "nova", "name": "Nova", "language": "zh-CN"},
		{"id": "shimmer", "name": "Shimmer", "language": "zh-CN"},
	}

	utils.SuccessResponse(c, gin.H{
		"voices": voices,
	})
}

// 批量转换接口
func (h *TTSHandler) BatchConvert(c *gin.Context) {
	var req struct {
		Texts []string `json:"texts" binding:"required"`
		model.TTSRequest
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[REQUEST] POST /api/v1/tts/batch | 客户端IP: %s | 请求体解析失败: %v", c.ClientIP(), err)
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	log.Printf("[REQUEST] POST /api/v1/tts/batch | 客户端IP: %s | 文本数量: %d | style: %q", c.ClientIP(), len(req.Texts), req.Style)

	if len(req.Texts) == 0 {
		log.Printf("[RESPONSE] POST /api/v1/tts/batch | 状态: 400 | 文本列表为空")
		utils.ErrorResponse(c, http.StatusBadRequest, "文本列表不能为空")
		return
	}

	if len(req.Texts) > 10 {
		log.Printf("[RESPONSE] POST /api/v1/tts/batch | 状态: 400 | 文本数量超过10")
		utils.ErrorResponse(c, http.StatusBadRequest, "一次最多转换10段文本")
		return
	}

	results, err := h.service.BatchConvert(req.Texts, &req.TTSRequest)
	if err != nil {
		log.Printf("[RESPONSE] POST /api/v1/tts/batch | 状态: 500 | 错误: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量转换失败: "+err.Error())
		return
	}

	log.Printf("[RESPONSE] POST /api/v1/tts/batch | 状态: 200 | 成功 | 转换数量: %d", len(results))
	utils.SuccessResponse(c, gin.H{
		"results": results,
	})
}