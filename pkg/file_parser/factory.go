package file_parser

import (
	"errors"

	"github.com/h2non/filetype"
)

// GetFileType 判断文件类型
//
// 传入文件前261字节的内容，通过Magic Bytes判断
//
// 使用 github.com/h2non/filetype 三方库
func GetFileType(buf []byte) (string, error) {
	// 检测类型
	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		// 未知类型，返回空字符串和错误
		return "", errors.New("Unknown Filetype.")
	}
	// 返回ext
	return kind.Extension, nil
}

// CheckTypeSupport 检查是否是支持的类型
func CheckTypeSupport(filetype string) bool {
	// 定义支持类型集合
	allowed := map[string]bool{
		"csv":  true,
		"docx": true,
		"pdf":  true,
		"txt":  true,
		"xlsx": true,
	}
	// 直接返回，不在集合中的会返回false
	return allowed[filetype]
}

// GetParser 文件解析器工厂函数：根据文件类型返回对应解析器
//
// 工厂模式
func GetParser(fileType string) FileParser {
	switch fileType {
	case "pdf":
		return &PDFParser{}
	case "xlsx":
		return &XLSXParser{}
	case "csv":
		return &CSVParser{}
	case "docx":
		return &DocxParser{}
	case "txt":
		return &TXTParser{}
	default:
		return nil
	}
}
