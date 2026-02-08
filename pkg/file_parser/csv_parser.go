package file_parser

import (
	"encoding/csv"
	"io"
	"strings"
)

// txt文件解析器
type CSVParser struct {
}

// Parse 解析csv文件内容
//
// 使用标准库
func (parser *CSVParser) Parse(file io.Reader) (string, error) {
	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 使用textBuilder进行拼接
	var textBuilder strings.Builder

	// 读取csv文件内容
	// 逐行读取
	for {
		// 这里的返回值是一个[]string,即每行记录的每个字段
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// 遍历单行记录的字段
		for col, field := range record {
			// 将字段添加到结果中
			textBuilder.WriteString(field)

			// 在字段间添加空格，如果是最后一个字段的话就不添加
			if col < len(record)-1 {
				textBuilder.WriteString(" ")
			}
		}

		// 在每行间添加一个逗号，保证后续chunk时，可以将记录切分完整
		textBuilder.WriteString(",")
	}

	// 获取最终提取的文本
	extractedText := textBuilder.String()

	return extractedText, nil
}
