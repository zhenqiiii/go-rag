package handlers

import (
	"go-rag/api/code"
	"go-rag/models"
	"go-rag/rag"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 文档操作的处理函数

// GetDocuments 获取文档列表
func GetDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从sqlite中读取文档id并返回
	}
}

// UploadDocument响应体
type UploadDocumentResponse struct {
	Success        bool     `json:"success"`
	DocumentID     string   `json:"document_id"`
	ChunkCount     int      `json:"chunk_count"`
	ChunkIDs       []string `json:"chunk_ids"`
	ProcessingTime int64    `json:"processing_time"` // 处理耗时（毫秒）
	Timestamp      string   `json:"timestamp"`       // 响应时间
}

// document upload路由处理
//
// 接收上传的文件
func UploadDocument(pipeline *rag.RAGPipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// 接收文件:文件通过form表单提交
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "获取文件失败：" + err.Error(),
				Code:    code.INVALID_REQUEST,
			})
			return
		}
		filename := file.Filename

		// 打开文件
		fileHandle, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "打开文件失败：" + err.Error(),
				Code:    code.INTERNAL_ERROR,
			})
			return
		}

		//读取
		fileBytes, err := io.ReadAll(fileHandle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "读取文件失败：" + err.Error(),
				Code:    code.INTERNAL_ERROR,
			})
			return
		}

		content := string(fileBytes)

		// 创建文档对象
		document := models.NewDocument(filename, content)

		// 调用pipeline存入向量数据库
		result, err := pipeline.IndexDocument(document)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "文档向量化失败：" + err.Error(),
				Code:    code.INTERNAL_ERROR,
			})
			return
		}

		// TODO：设置一个文档表-Mysql或者sqlite，暂时考虑只存放文档名称等必要信息，不存放内容

		// 返回响应
		c.JSON(http.StatusOK, UploadDocumentResponse{
			Success:        true,
			DocumentID:     result.DocumentID,
			ChunkCount:     result.ChunkCount,
			ChunkIDs:       result.ChunkIDs,
			ProcessingTime: time.Since(startTime).Milliseconds(),
			Timestamp:      time.Now().Format(time.RFC3339),
		})

	}
}

type DeleteDocumentResponse struct {
	Success       bool   `json:"success"`
	DocumentID    string `json:"document_id"`
	DeletedChunks int    `json:"deleted_chunks"` // 删除的切片数,这个数据应该在文档库中获取
	Timestamp     string `json:"timestamp"`
}

// DeleteDocument 删除文件路由
func DeleteDocument(pipeline *rag.RAGPipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取路径参数id
		documentID := c.Param("id")

		// TODO：检查文档是否存在（文档表中）

		// 在向量库中删除所有切片
		err := pipeline.DeleteDocument(documentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "删除向量失败：" + err.Error(),
				Code:    code.INTERNAL_ERROR,
			})
			return
		}

		// TODO：从文档表中删除

		// 返回响应
		c.JSON(http.StatusOK, DeleteDocumentResponse{
			Success:    true,
			DocumentID: documentID,
			Timestamp:  time.Now().Format(time.RFC3339),
		})

	}
}

// GetSpecificDocument 获取id对应的文档
func GetSpecificDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据id获取文档

	}
}
