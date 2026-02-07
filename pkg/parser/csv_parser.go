package parser

import "io"

// txt文件解析器
type CSVParser struct {
}

func (parser *CSVParser) Parse(file io.Reader) (string, error) {
	return "", nil
}
