package file_parser

import "io"

// txt文件解析器
type DocxParser struct {
}

func (parser *DocxParser) Parse(file io.Reader) (string, error) {
	return "", nil
}
