package rag

import (
	"fmt"
	"go-rag/models"
	"go-rag/pkg/embedding"
	"sync"
)

// 调用embedding模型api进行文档和query向量化

// Embedder 向量化组件接口
type Embedder interface {
	// EmbedChunks 批量向量化chunks
	//
	// 一般用于文档向量化
	EmbedChunks(chunks []models.Chunk) error

	//EmbedQuery 向量化查询文本
	EmbedQuery(query string) ([]float32, error)
}

// APIEmbedder 基于API的向量化器
//
// 拥有一个EmbeddingClient对象
type APIEmbedder struct {
	client    *embedding.EmbeddingClient
	batchSize int // 批量处理的大小
}

// NewAPIEmbedder 创建APIEmbedder向量化器
func NewAPIEmbedder(client *embedding.EmbeddingClient, batchSize int) *APIEmbedder {
	if batchSize <= 0 {
		batchSize = 10 // 默认批量处理的大小
	}
	return &APIEmbedder{
		client:    client,
		batchSize: batchSize,
	}
}

// EmbedChunks 批量向量化chunks
//
// 调用client的请求方法发送向量化请求，将响应返回到chunks
//
// 这里采用了并发批处理,然后chunkSize和embedding模型的上下文窗口限制以及API并发数限制都要协调一下
//
// TODO:添加并发数限制
func (e *APIEmbedder) EmbedChunks(chunks []models.Chunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// 提取文本
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	// 分批处理
	// （文档chunk后得到的chunks长度不固定，
	// 所以统一根据Embedder的batchSize来进行分批处理）
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	done := make(chan bool, 1)

	// 分批启动向量化线程
	for i := 0; i < len(texts); i += e.batchSize {
		end := i + e.batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]

		wg.Add(1)
		// 接收该批次的起始位置+文本内容切片
		go func(batchStart int, batchTexts []string) {
			defer wg.Done()

			// 调用向量化请求方法
			embeddings, err := e.client.Embed(batchTexts)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("批次 %d 向量化失败: %w", batchStart/e.batchSize, err):
				default:
				}
				return // 终止线程
			}

			// 向量化成功，将向量分配回chunks
			for j, emb := range embeddings {
				chunks[batchStart+j].Embedding = emb
			}
		}(i, batch)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan: // 监听错误
		return err
	case <-done: // 所有goroutine完成
		return nil
	}
}

// EmbedQuery 向量化查询文本
//
// 调用client的EmbedSingle方法
func (e *APIEmbedder) EmbedQuery(query string) ([]float32, error) {
	if query == "" {
		return nil, fmt.Errorf("查询文本不能为空")
	}

	vector, err := e.client.EmbedSingle(query)
	if err != nil {
		return nil, fmt.Errorf("向量化查询失败： %w", err)
	}
	return vector, nil
}
