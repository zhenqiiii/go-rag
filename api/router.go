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
		rag.POST("/query", handlers.SubmitQuery())

		// 文档管理:RESTful
		rag.GET("/documents", handlers.GetDocuments())           // 获取完整上传文档列表
		rag.GET("documents/:id", handlers.GetSpecificDocument()) //获取id对应的文档
		rag.POST("/documents", handlers.UploadDocument())        // 上传文档
		rag.DELETE("/documents/:id", handlers.DeleteDocument())  // 删除对应id文档

	}

	return r
}
