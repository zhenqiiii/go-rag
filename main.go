package main

import (
	"fmt"
	"go-rag/api"
	"go-rag/config"
	"go-rag/pkg/embedding"
	"go-rag/pkg/llm"
	"go-rag/rag"
	"log"
)

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败： %v", err)
	}

	// TODO:初始化数据库

	// 启动RAG服务
	// 根据config文件配置好pipeline
	// 然后注入到后端路由中供路由处理函数调用

	// 两个API客户端:embedding & llm
	embeddingClient := embedding.NewEmbeddingClient(
		config.EmbeddingAPI.BaseURL,
		config.EmbeddingAPI.APIKey,
		config.EmbeddingAPI.Model,
	)
	llmClient := llm.NewLLMClient(
		config.LLMAPI.BaseURL,
		config.LLMAPI.APIKey,
		config.LLMAPI.Model,
	)

	// 先创建store组件和embedder组件(retriever组件要复用，保证组件实例只有一个)
	store := rag.NewQdrantStore(
		config.Qdrant.Host,
		config.Qdrant.Port,
		config.Qdrant.Collection,
		config.Qdrant.Dimension,
	)
	// 初始化store组件（不过这样的话似乎pipeline的Init方法就多余了，之前没考虑到）
	if err := store.Init(); err != nil {
		log.Fatalf("初始化向量库失败：%v", err)
	}

	// embedder
	embedder := rag.NewAPIEmbedder(embeddingClient, 10)

	// 创建pipeline
	pipeline := rag.NewRAGPipeline(
		rag.NewTextChunker(config.RAG.ChunkSize, config.RAG.ChunkOverlap),
		embedder,
		store,
		rag.NewVectorRetriever(store, embedder, float32(config.RAG.ScoreThreshold)),
		rag.NewLLMGenerator(llmClient),
	)

	// 启动后端服务并注入rag服务和数据库（TODO）
	r := api.SetupRouter(pipeline)
	// 小马快跑 🐎💨
	r.Run(fmt.Sprintf(":%d", config.Server.Port))
}
