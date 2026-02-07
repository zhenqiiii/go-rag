package file_parser

import "io"

// FileParser接口定义
type FileParser interface {
	Parse(file io.Reader) (string, error) // 解析文件内容
}
