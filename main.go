package main

import (
	"go-rag/api"
	"go-rag/config"
	"log"
	"strconv"
)

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败： %v", err)
	}

	// 启动RAG服务
	// 根据config文件配置好pipeline
	// 然后注入到后端路由中供路由处理函数调用

	// 启动后端服务
	r := api.SetupRouter()
	// 小马快跑 🐎💨
	r.Run(strconv.Itoa(config.Server.Port))

	// // RAGPipeline整体测试
	// // chunker
	// chunker := rag.NewTextChunker(300, 20)

	// // embedder
	// client := embedding.NewEmbeddingClient(config.EmbeddingAPI.BaseURL, config.EmbeddingAPI.APIKey, config.EmbeddingAPI.Model)
	// embedder := rag.NewAPIEmbedder(client, 10)

	// // store
	// store := rag.NewQdrantStore(config.Qdrant.Host, config.Qdrant.Port, config.Qdrant.Collection, config.Qdrant.Dimension)

	// // retriever
	// retriever := rag.NewVectorRetriever(store, embedder, 0.5)

	// // generator
	// llmClient := llm.NewLLMClient(
	// 	config.LLMAPI.BaseURL,
	// 	config.LLMAPI.APIKey,
	// 	config.LLMAPI.Model,
	// )
	// generator := rag.NewLLMGenerator(llmClient)

	// // 创建流水线
	// testPipeline := rag.NewRAGPipeline(
	// 	chunker,
	// 	embedder,
	// 	store,
	// 	retriever,
	// 	generator,
	// )
	// // Init
	// testPipeline.Init()

	// // 读取AI.txt文件
	// content, err := os.ReadFile("AI.txt")
	// if err != nil {
	// 	log.Fatalf("文件读取错误: %v", err)
	// }
	// document := &models.Document{
	// 	ID:        "321",
	// 	Filename:  "AI",
	// 	Content:   string(content),
	// 	CreatedAt: time.Now(),
	// }

	// // 上传文件
	// IdxResult, err := testPipeline.IndexDocument(document)
	// if err != nil {
	// 	log.Fatalf("上传文件失败: %v", err)
	// }
	// fmt.Println(IdxResult)

	// 问答
	// query := "人工智能伦理的核心原则有哪些，每条都简单说一说"
	// response, err := testPipeline.Query(query, 10)
	// if err != nil {
	// 	log.Fatalf("回答生成失败: %v", err)
	// }

	// fmt.Printf("%+v", response)

}
