package file_parser

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/h2non/filetype"
)

// GetFileType 判断文件类型
//
// 传入文件前261字节的内容，通过Magic Bytes判断;txt文件通过后缀判断
//
// 使用 github.com/h2non/filetype 三方库
func GetFileType(buf []byte, filename string) (string, error) {
	// 检测类型
	kind, _ := filetype.Match(buf)
	if kind != filetype.Unknown {
		log.Println("=============================")
		log.Printf("文件类型: %s \n", kind.Extension)
		log.Println("=============================")
		return kind.Extension, nil

	}
	// 由于上面的三方库不支持txt(6...),txt单独判断
	// 通过后缀+内容判断
	extLower := strings.ToLower(filepath.Ext(filename))
	if extLower == ".txt" {
		// 进一步验证内容是否是txt
		if isPlainText(buf) {
			log.Println("=============================")
			log.Println("文件类型: txt ")
			log.Println("=============================")
			return "txt", nil
		}

	}

	// 确实是未知类型
	return "", errors.New("Unknown Filetype.")
}

// TODO:更加可靠地判断txt文件
// isPlainText 判断该文件是否为纯文本
func isPlainText(data []byte) bool {
	if len(data) == 0 {
		return true // 空文本视作合法
	}

	// 如果是UTF-8编码,直接通过
	// 但是有个疑问:传入的[]byte本身就是固定长度从文件中切割下来的,
	// 会不会导致末尾某个字符的byte没切完全,然后Valid方法不认为是utf-8?
	// 如果上面的设想属实,这个分支就永远不会执行
	if utf8.Valid(data) {
		return true
	}
	// 非 UTF-8：检查不可打印字符比例（保守策略）
	printable := 0
	for _, b := range data {
		// 允许：空格(32)~波浪号(126)，以及 \t(9), \n(10), \r(13)
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printable++
		}
	}
	return float64(printable)/float64(len(data)) > 0.90 // 至少 90% 可打印
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
