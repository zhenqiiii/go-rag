package models

// api相关的数据结构体

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Success bool   `json:"success"` // 一直都是false
	Error   string `json:"error"`   // 错误信息
	Code    string `json:"code"`    // 错误代码
}
