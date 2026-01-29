package rag

import (
	"fmt"
	"go-rag/models"
)

// 将所有RAG组件编排起来，组成一个完整RAG

// RAGPipeline RAG编排流水线
//
// 将所有RAG组件编排起来，组成完整RAG系统
//
// RAG的主要流程有三个： 1. 上传文件 2. 删除文件 3. 问答
type RAGPipeline struct {
	chunker   Chunker     //文档分块器
	embedder  Embedder    // 向量化器
	store     VectorStore // 向量存储
	retriever Retriever   // 检索器
	generator Generator   // 回答生成器
}

// NewRAGPipeline 创建RAG Pipeline
//
// 将所有组件传入（参数传入前进行设定）
//
// 由于上传文件时也需要向量化和存储，
// 所以pipeline中单独拥有embedder和store对象，便于逻辑区分
func NewRAGPipeline(
	chunker Chunker,
	embedder Embedder,
	store VectorStore,
	retriever Retriever,
	generator Generator,
) *RAGPipeline {
	return &RAGPipeline{
		chunker:   chunker,
		embedder:  embedder,
		store:     store,
		retriever: retriever,
		generator: generator,
	}
}

// IndexDocument 文档入库
//
// 将新文档添加到向量库中
//
// 流程：
// 1. Document -> Chunker -> Chunks : 文档分块
// 2. Chunks -> Embedder -> Chunks(已经向量化): 向量化
// 3. Chunks(带向量) -> Store(Qdrant): 存入向量库
//
// 返回:文档ID, 分块数量, 分块ID列表
func (p *RAGPipeline) IndexDocument(document *models.Document) (*models.IndexDocumentResponse, error) {
	// 文档分块
	chunks, err := p.chunker.Chunk(document)
	if err != nil {
		return nil, fmt.Errorf("文档分块失败: %w", err)
	}
	if len(chunks) == 0 {
		return nil, fmt.Errorf("文档分块结果为空")
	}

	// 分块向量化
	if err := p.embedder.EmbedChunks(chunks); err != nil {
		return nil, fmt.Errorf("向量化失败: %w", err)
	}

	// 存入向量库
	if err := p.store.AddChunks(chunks); err != nil {
		return nil, fmt.Errorf("存储向量失败: %w", err)
	}

	// 提取分块ID
	chunkIDs := make([]string, len(chunks))
	for i, chunk := range chunks {
		chunkIDs[i] = chunk.ID
	}

	// 返回入库结果
	response := &models.IndexDocumentResponse{
		DocumentID: document.ID,
		ChunkCount: len(chunks),
		ChunkIDs:   chunkIDs,
	}

	return response, nil
}

// Query RAG查询流程
//
// 根据用户问题,检索相关文档并生成回答
//
// 流程:
//  1. Query -> Embedder -> Vector: 问题向量化
//  2. Vector -> Retriever -> RetrievalResult : 检索相关文档
//  3. Retrieved chunks -> Generator : 生成回答
//
// 接收:问题和检索chunk数量topK
//
// 返回:回答内容,检索结果,使用的chunk数量
func (p *RAGPipeline) Query(query string, topK int) (*models.RAGResponse, error) {
	// 检索相关文档:retirever(embedder+store)
	retrievalResults, err := p.retriever.Retrieve(query, topK)
	if err != nil {
		return nil, fmt.Errorf("检索失败: %w", err)
	}

	// 没有检索到任何相关文档(分数太低,被阈值卡掉)
	if len(retrievalResults) == 0 {
		return &models.RAGResponse{
			Answer:     "无法根据知识库中的信息找到相关内容来回答你的问题.",
			Sources:    []models.RetrievalResult{},
			UsedChunks: 0,
		}, nil
	}

	// 从RetirevalResults中提取chunks
	contexts := make([]models.Chunk, len(retrievalResults))
	for i, result := range retrievalResults {
		contexts[i] = result.Chunk
	}

	// 将query和检索到的内容给到llm,生成回答
	fmt.Println("正在生成回答...")
	answer, err := p.generator.Generate(query, contexts)
	if err != nil {
		return nil, fmt.Errorf("生成回答失败: %w", err)
	}

	// 返回响应:RAGResponse(包含llm的回答)
	response := &models.RAGResponse{
		Answer:     answer,
		Sources:    retrievalResults, // 这里面包含了chunk和对应的相似度分数
		UsedChunks: len(contexts),    //chunk的个数(和retrievalResults的元素个数相同,这里只是用了contexts计算)
	}

	return response, nil
}

// DeleteDocument 删除文档
//
// 给定DocumentID,从知识库中删除对应文档的所有chunks
//
// 向量库和文档记录库都要删除(目前还没有写文档记录库就是)
func (p *RAGPipeline) DeleteDocument(documentID string) error {
	// 调用store组件的DeleteDocument方法
	if err := p.store.DeleteDocument(documentID); err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}

	return nil
}

// Init 初始化编排Pipeline
//
// 负责初始化向量库,启动时调用一次就行
func (p *RAGPipeline) Init() error {
	// 初始化向量存储
	if err := p.store.Init(); err != nil {
		return fmt.Errorf("初始化store组件(向量库)失败: %w", err)
	}

	return nil
}
