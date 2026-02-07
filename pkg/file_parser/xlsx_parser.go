package file_parser

import "io"

// xlsx文件解析器
type XLSXParser struct {
}

func (parser *XLSXParser) Parse(file io.Reader) (string, error) {
	return "", nil
}
