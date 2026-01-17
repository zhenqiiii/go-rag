package main

import (
	"fmt"
	"go-rag/config"
	"log"
)

func main() {
	// 加载配置文件
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败： %v", err)
	}

	fmt.Printf("配置调试：%+v", config)
}
