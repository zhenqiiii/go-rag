package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 文档操作的处理函数

// GetDocuments 获取文档列表
func GetDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从sqlite中读取文档id并返回
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

// DeleteDocument 删除文件路由
func DeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 接收文件信息
		// 进行删除操作：向量库删除所有切片
	}
}

// GetSpecificDocument 获取id对应的文档
func GetSpecificDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据id获取文档

	}
}
