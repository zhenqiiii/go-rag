package parser

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
