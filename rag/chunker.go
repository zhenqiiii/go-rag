package rag

import (
	"go-rag/models"
	"strings"
)

// 调用第三方的分片器对文档进行分片处理

// Chunker 文档分块器接口
//
// 供不同分块策略的实现
type Chunker interface {
	// Chunk 将文档分片
	//
	// 接收原始Document，根据策略切分并返回切分好的Chunk切片
	Chunk(document *models.Document) ([]models.Chunk, error)
}

// TextChunker 简单的基于语义的分块器
//
// 策略：以rune为基本单位进行切割，按照chunksize切分，
// 在最后一个标点符号出进行切分以优化，这里的一个rune可以看作一个字
type TextChunker struct {
	chunkSize    int // 分块大小（rune）
	chunkOverlap int // 分块重叠长度（rune）
}

// NewTextChunker创建分块器
func NewTextChunker(chunkSize, chunkOverlap int) *TextChunker {
	return &TextChunker{
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
	}
}

// Chunk 实现切分逻辑
//
// 这里按照rune切分，保证不会破坏utf-8编码的中文
func (tc *TextChunker) Chunk(document *models.Document) ([]models.Chunk, error) {
	// 清除space
	content := strings.TrimSpace(document.Content)       // 两头
	content = strings.Join(strings.Fields(content), " ") // 中间（先拆后拼）

	// 文档为空
	if len(content) == 0 {
		return nil, nil
	}

	// 转换为rune切片，确保按文字切分
	runes := []rune(content)
	// fmt.Print(len(runes))

	var chunks []models.Chunk
	start := 0 // 每个chunk的起始位置
	index := 0

	// 切分逻辑
	// 停止条件为start+tc.chunkOverlap >= 文档长度len(runes),可以理解为end已经到达文档末尾
	// 但这样似乎有点难理解,有待优化
	for start+tc.chunkOverlap < len(runes) {
		// 计算切片结束位置
		end := start + tc.chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		// 提取chunk,此时只按chunksize切分
		chunkRunes := runes[start:end]

		// 在标点符号处优化切分
		if end < len(runes) { // 保证后面还有内容，有必要切分
			// 在chunk末尾查找标点符号
			lastPunctPos := findLastPunctuation(chunkRunes)
			if lastPunctPos > 0 && lastPunctPos > len(chunkRunes)/2 {
				// 找到且位于chunk的后半部分，则进行切分
				chunkRunes = runes[start : start+lastPunctPos+1]
				end = start + len(chunkRunes) // end此时在最后一个标点符号处
			}
		}

		// 创建chunk
		chunk := models.NewChunk(document.ID, string(chunkRunes), index)
		chunks = append(chunks, *chunk)

		index++

		// 计算下一个chunk的起始位置(考虑overlap)
		start = end - tc.chunkOverlap
		if start < 0 {
			// chunkOverlap过大导致start重置,避免这种情况可以将重叠窗口设小一点
			// 推荐:   chunkSize: chunkOverlap = 10:1
			start = 0
		}

		// 避免死循环 (?)
		if start >= end {
			// 避免跳段(有内容没读到)
			start = end
		}
	}
	return chunks, nil

}

// findLastPunctuation 查找文本中最后一个标点符号的位置（rune计数）
//
// 接收转[]rune类型后的字符串,返回标点符号的索引（没找到返回-1）
func findLastPunctuation(runes []rune) int {
	// 定义标点符号
	punctuations := []rune{'。', '！', '？', '；', '.', '!', '?', ';'}

	// 从后往前遍历rune，找到最后一个标点符号
	for i := len(runes) - 1; i >= 0; i-- {
		if isPunctuation(runes[i], punctuations) {
			return i
		}
	}
	// 没找到标点符号则return-1
	return -1

}

// isPunctuation 判断字符是否是标点符号
//
// 将传入的字符和已定义的标点符号组比较
func isPunctuation(char rune, punctuations []rune) bool {
	for _, p := range punctuations {
		if char == p {
			return true
		}
	}
	return false
}
