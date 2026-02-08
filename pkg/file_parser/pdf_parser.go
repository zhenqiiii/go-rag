package file_parser

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// pdf文件解析器
type PDFParser struct{}

// Parse 从PDF文件中提取文本
//
// 接收文件的io.Reader接口，返回文本内容+error
func (parser *PDFParser) Parse(file io.Reader) (string, error) {
	// file 是一个io.Reader，只支持顺序读取，要转成io.ReadSeeker
	// 先进行二进制读取
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// 读入bytes.Reader(该reader实现了io.Seeker,可以作为io.ReadSeeker使用)
	byte_content := bytes.NewReader(content)

	// 创建pdf读取器
	pdfReader, err := model.NewPdfReader(byte_content)
	if err != nil {
		return "", err
	}

	// 获取pdf总页数,unipdf按页读取内容,所以后续遍历会用到
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	// 定义一个string.Builder拼接文本
	var textBuilder strings.Builder

	// 遍历所有页面,提取文本(页码从1开始)
	// 某一步出错直接返回错误,否则到最后检索的时候很难知道哪里出问题
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// 获取Page对象
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {

			return "", err
		}

		// 创建文本提取器
		textExtractor, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		// 提取页面文本
		text, err := textExtractor.ExtractText()
		if err != nil {
			return "", err
		}

		// 拼接到结果中
		// 在每页之间添加换行符,保持页面分隔
		textBuilder.WriteString(text)
		textBuilder.WriteString("\n\n")
	}

	// 获取最后得到的文本
	extractedText := textBuilder.String()

	// 检查是否成功提取文本
	if len(extractedText) == 0 {
		return "", errors.New("Fail to extract text from pdf.")
	}

	// 返回string类型文本内容
	return extractedText, nil
}
