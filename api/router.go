package api

import (
	"go-rag/api/handlers"
	"go-rag/rag"

	"github.com/gin-gonic/gin"
)

// 路由注册，设定路由引擎

// SetupRouter
//
// 传入配置好的pipeline，这样全局就只使用一个RAGPipeline实例了
func SetupRouter(pipeline *rag.RAGPipeline) *gin.Engine {
	r := gin.Default()
	// 中间件注册
	// TODO：加入一些必要的中间件把功能完善

	// RAG功能路由组
	rag := r.Group("api/rag")
	{
		// 查询路由
		rag.POST("/query", handlers.SubmitQuery(pipeline))

		// 文档管理:RESTful
		rag.GET("/documents", handlers.GetDocuments())            // 获取完整上传文档列表
		rag.GET("documents/:id", handlers.GetSpecificDocument())  //获取id对应的文档
		rag.POST("/documents", handlers.UploadDocument(pipeline)) // 上传文档
		rag.DELETE("/documents/:id", handlers.DeleteDocument())   // 删除对应id文档

	}

	return r
}
