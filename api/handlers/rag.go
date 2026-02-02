package handlers

import (
	"go-rag/api/code"
	"go-rag/models"
	"go-rag/rag"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// rag模块相关接口的处理函数

// 请求结构体
//
// 只给前端开放TopK的参数调整，相似度阈值由后端决定
type RAGQueryRequest struct {
	Query  string `json:"query" binding:"required"` // 用户问题，必填
	TopK   int    `json:"top_k"`                    // 检索返回文档数，可选，默认5
	Stream bool   `json:"stream"`                   // 流式输出开关，可选，默认false
}

// 请求响应体（非流式响应）
type RAGQueryResponse struct {
	Success    bool                     `json:"success"`     // 是否成功
	Answer     string                   `json:"answer"`      // 回答内容
	Sources    []models.RetrievalResult `json:"sources"`     // 检索到的文档chunks
	UsedChunks int                      `json:"used_chunks"` // 实际使用chunks数量
	Latency    int64                    `json:"latency"`     // 总耗时（ms）
	Timestamp  string                   `json:"timestamp"`   // 响应时间
}

// Query路由处理
func SubmitQuery(pipeline *rag.RAGPipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 解析请求
		var req RAGQueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "请求参数错误：" + err.Error(),
				Code:    code.INVALID_REQUEST,
			})
			return
		}

		// 设置TopK默认值
		if req.TopK <= 0 {
			req.TopK = 5
		}

		// 检查是否启动流式输出

		// stream = false,进行非流式处理
		handleNormalQuery(c, pipeline, req, startTime)

	}

}

// handleNormalQuery 处理非流式查询
func handleNormalQuery(c *gin.Context, pipeline *rag.RAGPipeline, req RAGQueryRequest, startTime time.Time) {

}
