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
		// 查询路由
		rag.POST("/query", handlers.RAGQuery())

		// 文档管理
		rag.POST("/upload", handlers.UploadDocument())
		rag.POST("/delete/:id", handlers.DeleteDocument())
		rag.GET("/documents", handlers.GetDocuments())
	}

	return r
}
