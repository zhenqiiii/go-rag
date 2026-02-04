package models

import (
	"time"

	"github.com/google/uuid"
)

// 文档分片
type Chunk struct {
	ID         string    `json:"id"`          // 唯一标识
	DocumentID string    `json:"document_id"` // 所属文档ID
	Content    string    `json:"content"`     // 分片内容
	Index      int       `json:"index"`       // 在文档中的索引(表示是该文档的第几个chunk,从0开始)
	Embedding  []float32 `json:"-"`           //	向量表示
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
}

// NewChunk 创建新的文档分片
func NewChunk(documentID string, content string, index int) *Chunk {
	return &Chunk{
		ID:         uuid.NewString(),
		DocumentID: documentID,
		Content:    content,
		Index:      index,
		Embedding:  nil, // Embed阶段进行填充
		CreatedAt:  time.Now(),
	}
}

// 检索结果
type RetrievalResult struct {
	Chunk Chunk   `json:"chunk"` // 检索到的分片
	Score float32 `json:"score"` // 相似度分数
}

// rag查询请求
type RAGRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// rag返回的响应，Answer由generator生成
type RAGResponse struct {
	Answer     string            `json:"answer"`      //回答
	Sources    []RetrievalResult `json:"sources"`     // 检索到的分块
	UsedChunks int               `json:"used_chunks"` // 实际使用的chunks数量
}

// 文档索引（切片入库）后返回的响应
type IndexDocumentResponse struct {
	DocumentID string   `json:"document_id"` // 原文档ID
	ChunkCount int      `json:"chunk_count"` //	分片数量
	ChunkIDs   []string `json:"chunk_ids"`   // 分片ID列表
}

// 删除文档请求
type DeleteDocumentRequest struct {
	DocumentID string `json:"document_id" binding:"required"`
}
