package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/model"
)

type TTSService struct {
	config *config.Config
	client *http.Client
}

func NewTTSService(cfg *config.Config) *TTSService {
	return &TTSService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// 调用 mimo-v2-tts 模型
func (s *TTSService) ConvertTextToSpeech(req *model.TTSRequest) (*model.TTSResponse, error) {
	// 构造文本内容，支持 style 标签
	content := req.Text
	if req.Style != "" {
		content = fmt.Sprintf("<style>%s</style>%s", req.Style, req.Text)
	}

	// 构造 chat completions 请求
	mimoReq := model.MimoChatRequest{
		Model: "mimo-v2-tts",
		Messages: []model.MimoMessage{
			{
				Role:    "assistant",
				Content: content,
			},
			{
				Role:    "assistant",
				Content: "太棒了！",
			},
		},
	}

	jsonData, err := json.Marshal(mimoReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	log.Printf("[MIMO-REQUEST] POST %s | model: %s | content长度: %d", s.config.MimoEndpoint, mimoReq.Model, len(content))

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", s.config.MimoEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.MimoAPIKey)

	// 发送请求
	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("[MIMO-ERROR] 调用失败: %v", err)
		return nil, fmt.Errorf("调用Mimo API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	log.Printf("[MIMO-RESPONSE] 状态码: %d | 响应体大小: %d bytes", resp.StatusCode, len(body))

	// 解析响应
	var mimoResp model.MimoChatResponse
	if err := json.Unmarshal(body, &mimoResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查错误
	if mimoResp.Error != nil {
		log.Printf("[MIMO-RESPONSE] API错误: %s", mimoResp.Error.Message)
		return &model.TTSResponse{
			Success: false,
			Message: mimoResp.Error.Message,
		}, nil
	}

	// 提取音频数据
	if len(mimoResp.Choices) == 0 || mimoResp.Choices[0].Message.Audio == nil {
		log.Printf("[MIMO-RESPONSE] 响应中无音频数据 | id: %s", mimoResp.ID)
		return &model.TTSResponse{
			Success: false,
			Message: "API未返回音频数据",
		}, nil
	}

	audioData, err := base64.StdEncoding.DecodeString(mimoResp.Choices[0].Message.Audio.Data)
	if err != nil {
		return nil, fmt.Errorf("解码音频数据失败: %v", err)
	}

	log.Printf("[MIMO-RESPONSE] 音频解码成功 | id: %s | 音频大小: %d bytes", mimoResp.ID, len(audioData))
	return &model.TTSResponse{
		Success:   true,
		Message:   "转换成功",
		AudioData: audioData,
		RequestID: mimoResp.ID,
	}, nil
}

// 批量文本转语音
func (s *TTSService) BatchConvert(texts []string, req *model.TTSRequest) ([]*model.TTSResponse, error) {
	results := make([]*model.TTSResponse, 0, len(texts))

	for _, text := range texts {
		req.Text = text
		result, err := s.ConvertTextToSpeech(req)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}