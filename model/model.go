package model

// TTS请求结构
type TTSRequest struct {
	Text  string `json:"text" binding:"required"`
	Voice string `json:"voice"` // 音色选择
	Style string `json:"style"` // 风格，如"东北话 开心"、"温柔"、"悄悄话"
}

// TTS响应结构
type TTSResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	AudioData []byte `json:"-"`        // 不在JSON中序列化
	RequestID string `json:"request_id,omitempty"`
}

// Mimo Chat Completions 请求结构
type MimoChatRequest struct {
	Model    string           `json:"model"`
	Messages []MimoMessage    `json:"messages"`
}

type MimoMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Mimo Chat Completions 响应结构
type MimoChatResponse struct {
	ID      string       `json:"id"`
	Choices []MimoChoice `json:"choices"`
	Error   *MimoError   `json:"error,omitempty"`
}

type MimoChoice struct {
	Message MimoChoiceMessage `json:"message"`
}

type MimoChoiceMessage struct {
	Role    string     `json:"role"`
	Content string     `json:"content"`
	Audio   *MimoAudio `json:"audio,omitempty"`
}

type MimoAudio struct {
	Data string `json:"data"` // base64编码的音频
}

type MimoError struct {
	Message string `json:"message"`
}