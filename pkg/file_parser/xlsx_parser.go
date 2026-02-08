package file_parser

import (
	"bytes"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

// xlsx文件解析器
type XLSXParser struct {
}

//突然感觉直接提取xlsx和csv文件的文本内容没什么用，需要换一种方式处理才能有效
// TODO: 用更为有效的方式处理csv & xlsx 数据

// Parse 解析xlsx文件内容
//
// 使用三方库 "github.com/xuri/excelize/v2"
func (parser *XLSXParser) Parse(file io.Reader) (string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 使用excelize从内存中打开,同样要转成*bytes.Reader
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	// 获取所有sheet名称
	sheetList := f.GetSheetList()

	// 只读取第一个工作表中的内容
	records, err := f.GetRows(sheetList[0])
	if err != nil {
		return "", err
	}

	// 使用textBuilder进行拼接
	var textBuilder strings.Builder

	// 遍历所有数据
	for row, record := range records {
		for col, field := range record {
			textBuilder.WriteString(field)
			// 使用空格分隔每个字段
			if col < len(record)-1 {
				textBuilder.WriteString(" ")
			}
		}
		// 每行后面加上逗号，方便后续chunk
		if row < len(records)-1 {
			textBuilder.WriteString(",")
		}
	}

	// 返回最终结果
	extractedContent := textBuilder.String()
	return extractedContent, nil

}
