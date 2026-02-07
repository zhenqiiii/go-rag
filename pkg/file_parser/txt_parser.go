package file_parser

import "io"

// txt文件解析器
type TXTParser struct {
}

// Parse 解析文件内容
func (parser *TXTParser) Parse(file io.Reader) (string, error) {
	// txt文件直接读取即可
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
