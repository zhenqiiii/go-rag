package api

import (
	"go-rag/api/handlers"

	"github.com/gin-gonic/gin"
)

// 路由注册，设定路由引擎

// SetupRouter
func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 中间件注册

	// RAG功能路由组
	rag := r.Group("api/rag")
	{
		rag.POST("/query", handlers.RAGQuery)
		rag.POST("/upload", handlers.UploadDocument)
	}

	return r
}
