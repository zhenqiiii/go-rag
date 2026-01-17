package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// rag模块相关接口的处理函数

// Query路由处理
func RAGQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取问题
		// 交给RAG
		// 返回回答
		c.JSON(http.StatusOK, gin.H{
			"msg": "RAGQuery",
		})
	}

}

// document upload路由处理
func UploadDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 接收文件

		// 交给RAG向量化后存入向量库
		c.JSON(http.StatusOK, gin.H{
			"msg": "UploadDocument",
		})
	}
}
