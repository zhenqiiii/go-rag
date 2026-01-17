package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 程序配置
type Config struct {
	// 后端配置
	Server ServerConfig `mapstructure:"server"`

	// Qdrant配置
	Qdrant QdrantConfig `mapstructure:"qdrant"`

	// Embedding API配置
	EmbeddingAPI EmbeddingAPIConfig `mapstructure:"embedding_api"`

	// LLM对话模型配置
	LLMAPI LLMAPIConfig `mapstructure:"llm_api"`

	// RAG参数
	RAG RAGConfig `mapstructure:"rag"`
}

// web后端服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// Qdrant向量数据库配置
type QdrantConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Collection string `mapstructure:"collection"`
}

// Embedding API配置
type EmbeddingAPIConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	APIKey    string `mapstructure:"api_key"`
	Model     string `mapstructure:"model"`
	Dimension int    `mapstructure:"dimension"` //使用BGE-m3的话，发请求时这个字段用不上
}

// LLM API配置
type LLMAPIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
}

// RAG系统参数配置
type RAGConfig struct {
	ChunkSize      int     `mapstructure:"chunk_size"`      // 分块大小（字符数）
	ChunkOverlap   int     `mapstructure:"chunk_overlap"`   // 分块重叠字符数
	TopK           int     `mapstructure:"top_k"`           // 检索返回的top数量
	ScoreThreshold float64 `mapstructure:"score_threshold"` // 相似度阈值
}

// LoadConfig 从文件中加载配置
//
// configpath: config/config.yaml
func LoadConfig(configpath string) (*Config, error) {
	// 使用传入的路径
	// SetConfigFile方法接收文件的绝对路径,
	// 因为是在main中调用，所以使用config.yaml时要传入config/config.yaml
	viper.SetConfigFile(configpath)
	viper.SetConfigType("yaml")

	// 设置默认值	优先级:配置文件> 默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败： %w", err)
	}

	return &config, nil

}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)

	// Qdrant默认配置
	viper.SetDefault("qdrant.host", "localhost")
	viper.SetDefault("qdrant.port", 6333)
	viper.SetDefault("qdrant.collection", "documents")

	// RAG默认参数
	viper.SetDefault("rag.chunk_size", 500)
	viper.SetDefault("rag.chunk_overlap", 50)
	viper.SetDefault("rag.top_k", 5)
	viper.SetDefault("rag.score_threshold", 0.7)
}

// validateConfig 验证配置有效性
func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", config.Server.Port)
	}

	// 验证Qdrant配置
	if config.Qdrant.Host == "" {
		return fmt.Errorf("Qdrant host不能为空")
	}

	// 验证Embedding API配置
	if config.EmbeddingAPI.APIKey == "" {
		return fmt.Errorf("Embedding API Key不能为空")
	}
	if config.EmbeddingAPI.BaseURL == "" {
		return fmt.Errorf("Embedding Base URL不能为空")
	}
	if config.EmbeddingAPI.Model == "" {
		return fmt.Errorf("Embedding Model不能为空")
	}
	if config.EmbeddingAPI.Dimension <= 0 {
		return fmt.Errorf("Embedding Dimension必须大于0")
	}

	// 验证LLM API配置
	if config.LLMAPI.APIKey == "" {
		return fmt.Errorf("LLM API Key不能为空")
	}
	if config.LLMAPI.BaseURL == "" {
		return fmt.Errorf("LLM Base URL不能为空")
	}
	if config.LLMAPI.Model == "" {
		return fmt.Errorf("LLM Model不能为空")
	}

	// 验证RAG参数
	if config.RAG.ChunkSize <= 0 {
		return fmt.Errorf("Chunk Size必须大于0")
	}
	if config.RAG.ChunkOverlap < 0 {
		return fmt.Errorf("Chunk Overlap不能为负数")
	}
	if config.RAG.TopK <= 0 {
		return fmt.Errorf("Top K必须大于0")
	}
	if config.RAG.ScoreThreshold < 0 || config.RAG.ScoreThreshold > 1 {
		return fmt.Errorf("Score Threshold必须在0-1之间")
	}

	return nil

}
