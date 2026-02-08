package file_parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// txt文件解析器
type DocxParser struct {
}

// TextElement 用于存放从xml文件中解析出来的文本
type TextElement struct {
	Body string `xml:",chardata"`
}

// Parse 解析Docx文件内容
//
// 基于zip & xml标准库实现docx纯文本提取
//
// docx文件本质是一个zip文件,内部的word/document.xml包含了其文本
//
// 所以文本提取思路为: 解压docx对应的zip文件,然后解析其中的xml文件
//
// 这种方式只适用于简单文档
func (parser *DocxParser) Parse(file io.Reader) (string, error) {
	// 先转化到io.ReadSeeker
	data, ok := file.(*bytes.Reader)
	if !ok {
		return "", errors.New("errors while transfering io.Reader to bytes.Reader")
	}

	// 创建解压器
	zipReader, err := zip.NewReader(data, int64(data.Len()))
	if err != nil {
		return "", err
	}

	// 获取document.xml文件
	var docXML []byte
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			// 进行常规提取即可
			content, err := f.Open()
			if err != nil {
				return "", err
			}
			docXML, _ = io.ReadAll(content)
			content.Close()
			break
		}
	}

	// 如果docXML为空,可能是没在解压得到的文件中找到word/document.xml
	if docXML == nil {
		return "", errors.New("document.xml not found")
	}

	// 提取xml文件中的所有文本节点:使用xml标准库
	var buf strings.Builder
	decoder := xml.NewDecoder(bytes.NewReader(docXML))
	for {
		// Token()方法每次返回下一个token
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		// 判断该token是否是StartElement,即xml中的开始标签
		// 是则进行判断是否为<w:t>标签,xml文件中标签<w:t>(其标签名为t)包裹的内容就是文本
		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Local == "t" {
				// 存入TextElement结构体
				var t TextElement
				decoder.DecodeElement(&t, &se)
				// 写入结果
				buf.WriteString(t.Body)
			}
		}
	}
	return buf.String(), nil
}
