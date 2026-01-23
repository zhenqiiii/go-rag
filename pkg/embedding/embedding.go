package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 调用embedding模型逻辑

// EmbeddingClient EmbeddingAPI的客户端
type EmbeddingClient struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// NewEmbeddingClient 创建Embedding客户端
func NewEmbeddingClient(baseURL, apiKey, model string) *EmbeddingClient {
	return &EmbeddingClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{},
	}
}

// EmbeddingRequest Embedding API请求体
//
// 即将请求payload中的几个字段组合成一个结构体
type EmbeddingRequest struct {
	Input []string `json:"input"` // 使用string切片支持一个请求中处理多个input
	Model string   `json:"model"`
}

// EmbeddingResponse Embedding API响应体
type EmbeddingResponse struct {
	Object string `json:"object"` // 这个是[]string还是string？看硅基流动是一个string切片，但里面始终是一个元素
	Model  string `json:"model"`
	Data   []struct {
		Object    string    `json:"Object"`    // 任务类型（可以这么翻译吧应该）
		Embedding []float32 `json:"embedding"` // 向量
		Index     int       `json:"index"`     // 索引,即模型API返回的向量的顺序
	} `json:"data"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Embed方法 发送embed请求，批量向量化文本
//
// 接收要embed的string切片，返回[]float32类型向量的切片
func (c *EmbeddingClient) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("输入文本不能为空")
	}

	// 构建请求体
	reqBody := EmbeddingRequest{
		Input: texts,
		Model: c.model,
	}

	// json序列化
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/embeddings", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败：%w", err)
	}
	defer resp.Body.Close()

	//读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败：%w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API请求失败, http状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	// 200，向量化成功，解析响应到响应结构体
	var result EmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败： %w", err)
	}

	// 提取向量
	embeddings := make([][]float32, len(result.Data))
	for _, item := range result.Data {
		embeddings[item.Index] = item.Embedding
	}
	// fmt.Print(len(embeddings[0]))
	return embeddings, nil
}

// EmbedSingle 向量化单个文本
//
// 其实就是封装了Embed方法，但为了区分两种情况还是写了一下
func (c *EmbeddingClient) EmbedSingle(text string) ([]float32, error) {
	embeddings, err := c.Embed([]string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("未返回向量")
	}

	return embeddings[0], nil
}
