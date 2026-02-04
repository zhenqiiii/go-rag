package models

import (
	"time"

	"github.com/google/uuid"
)

// 原始文档
type Document struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewDocument 创建新文档
func NewDocument(filename, content string) *Document {
	return &Document{
		ID:        uuid.New().String(),
		Filename:  filename,
		Content:   content,
		CreatedAt: time.Now(),
	}
}
