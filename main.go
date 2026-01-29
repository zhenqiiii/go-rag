package main

import (
	"fmt"
	"go-rag/config"
	"go-rag/models"
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

	// chunker := rag.NewTextChunker(200, 20)
	// chunks, err := chunker.Chunk(document)
	// if err != nil {
	// 	log.Fatalf("切割错误: %v", err)
	// }

	// APIEmbedder测试
	client := embedding.NewEmbeddingClient(config.EmbeddingAPI.BaseURL, config.EmbeddingAPI.APIKey, config.EmbeddingAPI.Model)
	embedder := rag.NewAPIEmbedder(client, 10)
	// embedder.EmbedChunks(chunks) //向量化并将向量返回到chunks中

	// store 测试
	store := rag.NewQdrantStore(config.Qdrant.Host, config.Qdrant.Port, config.Qdrant.Collection, config.Qdrant.Dimension)
	// c初始化一下
	err = store.Init()
	if err != nil {
		log.Fatalf("初始化store组件实例失败: %v", err)
	}
	// // 添加chunk
	// err = store.AddChunks(chunks)
	// if err != nil {
	// 	log.Fatalf("添加chunk失败: %v", err)
	// }
	// // 向量化query并进行search
	query := "Golang的应用场景有哪些？"
	embeddedQuery, err := embedder.EmbedQuery(query)
	if err != nil {
		log.Fatalf("向量化query失败: %v", err)
	}
	results, err := store.Search(embeddedQuery, 3)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	// 拿到chunks
	contextText := make([]models.Chunk, len(results))
	for _, result := range results {
		contextText = append(contextText, result.Chunk)
	}

	llmClient := llm.NewLLMClient(
		config.LLMAPI.BaseURL,
		config.LLMAPI.APIKey,
		config.LLMAPI.Model,
	)
	generator := rag.NewLLMGenerator(llmClient)

	answer, err := generator.Generate(query, contextText)
	if err != nil {
		log.Fatalf("获取回答失败: %v", err)
	}

	fmt.Print(answer)

}
