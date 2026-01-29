package main

import (
	"fmt"
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

	// fmt.Printf("配置调试：%+v", config)

	// // TextChunker测试
	// content, err := os.ReadFile("Golang.txt")
	// if err != nil {
	// 	log.Fatalf("文件读取错误: %v", err)
	// }
	// // fmt.Print(string(content))
	// document := &models.Document{
	// 	ID:        "123",
	// 	Filename:  "Golang",
	// 	Content:   string(content),
	// 	CreatedAt: time.Now(),
	// }

	// RAGPipeline整体测试
	// chunker
	chunker := rag.NewTextChunker(300, 20)

	// embedder
	client := embedding.NewEmbeddingClient(config.EmbeddingAPI.BaseURL, config.EmbeddingAPI.APIKey, config.EmbeddingAPI.Model)
	embedder := rag.NewAPIEmbedder(client, 10)

	// store
	store := rag.NewQdrantStore(config.Qdrant.Host, config.Qdrant.Port, config.Qdrant.Collection, config.Qdrant.Dimension)

	// retriever
	retriever := rag.NewVectorRetriever(store, embedder, 0.5)

	// generator
	llmClient := llm.NewLLMClient(
		config.LLMAPI.BaseURL,
		config.LLMAPI.APIKey,
		config.LLMAPI.Model,
	)
	generator := rag.NewLLMGenerator(llmClient)

	// 创建流水线
	testPipeline := rag.NewRAGPipeline(
		chunker,
		embedder,
		store,
		retriever,
		generator,
	)
	// Init
	testPipeline.Init()

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
	query := "人工智能伦理的核心原则有哪些，每条都简单说一说"
	response, err := testPipeline.Query(query, 10)
	if err != nil {
		log.Fatalf("回答生成失败: %v", err)
	}

	fmt.Printf("%+v", response)

	// // // 向量化query并进行search
	// query := "Golang的应用场景有哪些？"
	// embeddedQuery, err := embedder.EmbedQuery(query)
	// if err != nil {
	// 	log.Fatalf("向量化query失败: %v", err)
	// }
	// results, err := store.Search(embeddedQuery, 3)
	// if err != nil {
	// 	log.Fatalf("查询失败: %v", err)
	// }

	// // 拿到chunks
	// contextText := make([]models.Chunk, len(results))
	// for _, result := range results {
	// 	contextText = append(contextText, result.Chunk)
	// }

	// answer, err := generator.Generate(query, contextText)
	// if err != nil {
	// 	log.Fatalf("获取回答失败: %v", err)
	// }

	// fmt.Print(answer)

}
