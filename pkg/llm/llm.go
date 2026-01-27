package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// llm客户端

// LLMClient LLM客户端
//
// 负责调用API进行对话生成
type LLMClient struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// NewLLMClient 创建LLM客户端
//
// 参数(从config传入)：
// - baseURL：api url
// - apiKey ：密钥
// - model： 模型名称
func NewLLMClient(baseURL, apiKey, model string) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{},
	}
}

// Message 对话消息结构
//
// 定义一个Message切片存放对话历史
//
// Q: 如何实现上下文？
type Message struct {
	Role    string `json:"role"`    // 一般有system/user/assistant
	Content string `json:"content"` // 内容
}

// ChatRequest 聊天请求体
//
// 有多个定制化字段，这里先上几个必须和常用的
type ChatRequest struct {
	Model       string    `json:"model"`                 // 模型
	Messages    []Message `json:"messages"`              // chat历史
	Temperature float64   `josn:"temperature,omitempty"` // 温度参数，控制随机性
	MaxTokens   int       `json:"max_tokens,omitempty"`  // 生成的最大token数，输入token+MaxTokens不能超过模型上下文
	Stream      bool      `json:"stream,omitempty"`      // 流式输出
}

// ChatResponse chat响应体
type ChatResponse struct {
	ID      string `json:"id"`      // 响应ID
	Object  string `json:"object"`  // 类型
	Created int    `json:"created"` // 创建时间
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"` // 结束原因
	} `json:"choices"`
	Usage struct { // 这个结构体字段贼多，挑了三个
		PromptTokens     int `json:"prompt_tokens"`     // 输入token数
		CompletionTokens int `json:"completion_tokens"` // 输出token数
		TotalTokens      int `json:"total_tokens"`      //总token数，是上面两个字段的总和，不能超过模型上下文
	} `json:"usage"`
}

// Chat 拼接prompt,发送聊天请求
//
// 参数:
// - systemPrompt :系统提示词(定义角色和行为)
// - userPrompt: 用户输入Query
// - temperature: 温度
// -maxTokens : 最大生成token数,负责限制输出
//
// 返回:LLM的回答文本
func (c *LLMClient) Chat(systemPrompt, userPrompt string, temperature float64, maxTokens int) (string, error) {
	// 构建对话历史
	// systemPrompt限定AI角色
	// userPrompt包含用户问题和retrieve得到的文档内容
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// 构建请求体
	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      false,
	}

	// 序列化为json
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建请求
	// 完整URL: baseURL + "/chat/completions"
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer((jsonData)))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API请求失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应到ChatResponse
	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取回答
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API未返回任何回答")
	}

	return result.Choices[0].Message.Content, nil

}
