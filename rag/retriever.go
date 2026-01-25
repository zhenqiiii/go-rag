package rag

import (
	"fmt"
	"go-rag/models"
)

// 向量检索组件

// Retriever向量检索器组件接口
type Retriever interface {
	// Retrieve 根据query执行检索操作
	//
	// 参数:
	// - query : 用户查询的文本
	//  -topK: 返回的结果数量
	// 返回: retrievalResult切片(包含chunk和相似度分数)
	Retrieve(query string, topK int) ([]models.RetrievalResult, error)
}

// VectorRetriever 检索器实现
//
// 封装embedder store 组件: EmbedQuery() -> Search()
type VectorRetriever struct {
	embedder       Embedder    // 向量化器,将query转换为向量
	store          VectorStore // 向量存储,执行搜索操作
	scoreThreshold float32     // 相似度阈值,用于从搜索结果中过滤低质量结果
}

// NewVectorRetriever 创建向量检索器
//
// 参数:
// - store: 向量存储组件
// - embedder: 向量化器
// - scoreThreshold: 相似度阈值,默认0.7(在config中设置)
func NewVectorRetriever(store VectorStore, embedder Embedder, scoreThreshold float32) *VectorRetriever {
	return &VectorRetriever{
		store:          store,
		embedder:       embedder,
		scoreThreshold: scoreThreshold,
	}
}

// Retrieve 执行检索
//
// 文本向量化 -> 进行搜索得到topK个结果 -> 过滤低于阈值的结果 -> 返回
func (vr *VectorRetriever) Retrieve(query string, topK int) ([]models.RetrievalResult, error) {
	// 向量化query:调用embedder的EmbedQuery方法
	queryVector, err := vr.embedder.EmbedQuery(query)
	if err != nil {
		return nil, fmt.Errorf("向量化Query失败: %w", err)
	}

	// 进行查询:调用store的Search方法（topK传入）
	results, err := vr.store.Search(queryVector, topK)
	if err != nil {
		return nil, fmt.Errorf("向量搜索失败: %w", err)
	}

	// 根据scoreThreshhold过滤部分低质量结果
	// 返回的相似度分数应该在-1~1之间，越接近1越相似
	var filteredResults []models.RetrievalResult
	for _, result := range results {
		if result.Score >= vr.scoreThreshold {
			filteredResults = append(filteredResults, result)
		}
	}

	// 过滤后没有结果，返回空slice
	if len(filteredResults) == 0 {
		return []models.RetrievalResult{}, nil
	}

	return filteredResults, nil
}
