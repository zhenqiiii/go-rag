package main

import (
	"fmt"
	"go-rag/config"
	"go-rag/models"
	"go-rag/pkg/embedding"
	"go-rag/rag"
	"log"
	"os"
	"time"
)

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败： %v", err)
	}

	// fmt.Printf("配置调试：%+v", config)

	// TextChunker测试
	content, err := os.ReadFile("Golang.txt")
	if err != nil {
		log.Fatalf("文件读取错误: %v", err)
	}
	// fmt.Print(string(content))
	document := &models.Document{
		ID:        "123",
		Filename:  "Golang",
		Content:   string(content),
		CreatedAt: time.Now(),
	}

	chunker := rag.NewTextChunker(200, 20)
	chunks, err := chunker.Chunk(document)
	if err != nil {
		log.Fatalf("切割错误: %v", err)
	}

	// APIEmbedder测试
	client := embedding.NewEmbeddingClient(config.EmbeddingAPI.BaseURL, config.EmbeddingAPI.APIKey, config.EmbeddingAPI.Model)
	embedder := rag.NewAPIEmbedder(client, 10)
	embedder.EmbedChunks(chunks)

	for _, chunk := range chunks {
		fmt.Printf("\n %+v \n", chunk)
	}
}
