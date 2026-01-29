package rag

import (
	"fmt"
	"go-rag/models"
	"go-rag/pkg/llm"
	"strings"
)

// 回答内容生成

// Generator 回答生成器接口
//
// 根据检索到的chunk，生成回答
type Generator interface {
	// Generate 生成回答
	//
	// 参数：
	// - query：用户问题
	// - contexts： 根据query检索到的chunks
	// 返回：回答文本
	Generate(query string, contexts []models.Chunk) (string, error)
}

// LLMGenerator 生成器接口实现
type LLMGenerator struct {
	llmClient *llm.LLMClient // LLM API客户端
}

// NewLLMGenerator 创建LLMGenerator
//
// 参数：
// - llmClient: 已经初始化的LLM客户端
func NewLLMGenerator(llmClient *llm.LLMClient) *LLMGenerator {
	return &LLMGenerator{
		llmClient: llmClient,
	}
}

// SystemPrompt 系统提示词
//
// 定义AI角色和行为
const SystemPrompt = `你是一个智能问答助手。请根据以下参考信息回答用户的问题。

【重要规则】：
1. 只能基于提供的参考信息回答问题
2. 如果参考信息中没有相关内容，请诚实告知"根据提供的信息无法回答"
3. 回答时要引用参考信息中的具体内容
4. 不要编造信息或使用你的外部知识
5. 回答要简洁、准确、有条理
6. 如果参考信息有矛盾，请指出矛盾之处

参考信息如下：`

// Generate 使用默认提示词生成回答
//
// 执行流程:
// 1. 构建包含上下文的用户提示词
// 2. 调用大模型生成回答
// 3. 返回回答
func (lg *LLMGenerator) Generate(query string, contexts []models.Chunk) (string, error) {
	// 拼接上下文:调用buildContextText方法
	contextText := lg.buildContextText(contexts)

	// 构建用户提示词:query+contextText（检索到的参考信息）
	userPrompt := lg.buildUserPrompt(query, contextText)

	// 调用llm生成回答
	// Temperature = 0.3	,基于事实回答
	// MaxTokens = 1024, 这个参数限制模型输出的长度，不包括输入
	answer, err := lg.llmClient.Chat(
		SystemPrompt,
		userPrompt,
		0.3,
		1024,
	)
	if err != nil {
		return "", fmt.Errorf("LLM生成失败： %w", err)
	}

	return answer, nil

}

// buildContextText 构建上下文文本
//
// 将多个chunk拼接起来,并对每个chunk使用序号标记
func (lg *LLMGenerator) buildContextText(contexts []models.Chunk) string {
	if len(contexts) == 0 {
		return "(无相关信息)"
	}

	var sb strings.Builder

	// 遍历所有传入的chunks,用序号进行标记(从1开始)
	for i, chunk := range contexts {
		sb.WriteString(fmt.Sprintf("\n[参考%d] %s", i+1, chunk.Content))
	}

	return sb.String()
}

// buildUserPrompt 构建用户提示词
//
// 将处理好的参考信息和用户query组成userPrompt，供llm读取
func (lg *LLMGenerator) buildUserPrompt(query string, contextText string) string {
	// 三段式： 1. 参考信息 2. 用户query 3. 可选的要求
	return fmt.Sprintf(`参考信息：
%s

用户问题：
%s

请根据上述参考信息回答用户问题，并在回答中注明信息来源。`,
		contextText,
		query,
	)
}

// GenerateWithPrompt 使用自定义的系统提示词生成回答
//
// 使用传入的提示词
func (lg *LLMGenerator) GenerateWithPrompt(query string, contexts []models.Chunk, systemPrompt string) (string, error) {
	// 同样是处理一下参考文档
	contextText := lg.buildContextText(contexts)

	// 构建用户提示词
	userPrompt := lg.buildUserPrompt(query, contextText)

	// 调用对话模型
	answer, err := lg.llmClient.Chat(
		systemPrompt, // 使用传入的自定义系统提示词
		userPrompt,
		0.3,
		500,
	)

	if err != nil {
		return "", fmt.Errorf("LLM生成失败: %w", err)
	}

	return answer, nil
}
