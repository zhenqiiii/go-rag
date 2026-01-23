package rag

import (
	"context"
	"fmt"
	"go-rag/models"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

// 向量存储

// VectorStore 向量存储组件接口
//
// 基本操作：添加、搜索、删除、初始化
type VectorStore interface {
	// AddChunks 批量添加chunk到向量库
	//
	// 每个chunk存作一个点，包含向量和payload（MetaData）
	AddChunks(chunks []models.Chunk) error

	// Search 根据query向量搜索最相似的chunks
	//
	// 返回top-K个最相似向量点，按相似度降序排列
	Search(queryVector []float32, topK int) ([]models.RetrievalResult, error)

	// DeleteDocument 删除指定文档的所有chunks
	//
	// 根据document_id删除所有相关向量point
	DeleteDocument(documentID string) error

	// Init 初始化向量库
	//
	// 若collection不存在，创建并配置向量维度和距离度量方式
	Init() error
}

// QdrantStore 基于Qdrant实现的向量存储库组件
type QdrantStore struct {
	client     *qdrant.Client // Qdrant GoSDK客户端
	collection string         // 集合名称,目前采用多租户形式，只用一个集合
	dimension  int            // 向量维度:1024
}

// NewQdrantStore 创建QdrantStore组件
//
// 参数：
//   - host: Qdrant服务地址
//   - port: Qdrant服务端口
//   - collection: 集合名称
//   - dimension: 向量维度
func NewQdrantStore(host string, port int, collection string, dimension int) *QdrantStore {
	// 创建Qdrant客户端连接
	client, _ := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})

	return &QdrantStore{
		client:     client,
		collection: collection,
		dimension:  dimension,
	}
}

// Init 初始化Qdrant存储
//
// 检查collection是否存在,不存在则创建
func (qs *QdrantStore) Init() error {
	ctx := context.Background()

	// 获取collection列表
	collections, err := qs.client.ListCollections(ctx)
	if err != nil {
		return fmt.Errorf("获取collection列表失败: %w", err)
	}
	// 检查目标collection是否存在
	var collectionExists bool
	for _, col := range collections {
		if col == qs.collection {
			collectionExists = true
			break
		}
	}
	// 不存在则创建
	if !collectionExists {
		log.Printf("创建collection: %s, 向量维度: %d", qs.collection, qs.dimension)

		// 配置向量参数
		vectorParams := qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(qs.dimension),   // 指定向量维度
			Distance: qdrant.Distance_Cosine, // 采用余弦相似度作为度量指标
		})

		// 创建collection
		err = qs.client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: qs.collection,
			VectorsConfig:  vectorParams,
			// 后面可以配置一下HNSW索引,目前还不知道这个东西有什么用
		})

		if err != nil {
			return fmt.Errorf("创建collection失败: %w", err)
		}
		log.Printf("Collection %s 创建成功", qs.collection)
	} else {
		log.Printf("Collection %s 已存在", qs.collection)
	}

	return nil
}

// AddChunks 批量添加chunks到Qdrant
//
// 每个chunk作为一个Point:包含
// - id: 点的唯一标识,使用chunk_ID
// - vector: 向量, 使用chunk.embedding
// - payload : 元数据,包含document_id, chunk_index, content
func (qs *QdrantStore) AddChunks(chunks []models.Chunk) error {
	ctx := context.Background()

	// 点切片points
	points := make([]*qdrant.PointStruct, 0, len(chunks))
	// 从chunks中批量导入数据到points
	for _, chunk := range chunks {
		// 创建payload:document_id, chunk_index, content
		payload := map[string]interface{}{
			"document_id": chunk.DocumentID, // 所属文档ID
			"chunk_index": chunk.Index,      // chunk在文档中的位置
			"content":     chunk.Content,    // 文本内容
		}

		// 创建点结构
		point := &qdrant.PointStruct{
			Id:      qdrant.NewID(chunk.ID),
			Vectors: qdrant.NewVectorsDense(chunk.Embedding),
			Payload: qdrant.NewValueMap(payload),
		}
		points = append(points, point)
	}

	// 插入点
	_, err := qs.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: qs.collection,
		Points:         points,
	})
	if err != nil {
		return fmt.Errorf("批量插入向量失败: %w", err)
	}
	log.Printf("成功插入 %d 个chunks到Qdrant", len(chunks))
	return nil
}

// Search 搜索相似的向量
//
// 根据query向量,在向量库中搜索top-k个最相似chunk
//
// 返回chunk+相似度分数
func (qs *QdrantStore) Search(queryVector []float32, topK int) ([]models.RetrievalResult, error) {
	ctx := context.Background()

	// 进行向量搜索
	searchResult, err := qs.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: qs.collection,
		Query:          qdrant.NewQueryDense(queryVector),
		Limit:          qdrant.PtrOf(uint64(topK)),
		WithPayload:    qdrant.NewWithPayload(true),
		ScoreThreshold: nil,
		// 查看这个QueryPoints结构体的字段了解更多可以设置的参数
		// Params,Fliter等等
		// 目前的策略是简单的基于余弦相似度的最近邻搜索,返回top-K个分数最高的chunk,分数阈值后面看情况设
	})
	if err != nil {
		return nil, fmt.Errorf("向量搜索失败: %w", err)
	}

	// 转换搜索结果到定义的RetrivalResult结构体中
	results := make([]models.RetrievalResult, len(searchResult))
	for i, point := range searchResult {
		// 提取元数据
		documentID := point.Payload["document_id"].GetStringValue()
		chunkIndex := point.Payload["chunk_index"].GetIntegerValue()
		content := point.Payload["content"].GetStringValue()

		// 创建chunk对象
		chunk := models.Chunk{
			ID:         point.Id.GetUuid(),
			DocumentID: documentID,
			Content:    content,
			Index:      int(chunkIndex),
		}

		// 创建检索结果
		results[i] = models.RetrievalResult{
			Chunk: chunk,
			Score: point.Score, // 相似度分数(余弦相似度,-1~1,越大越相似)
		}
	}

	return results, nil
}

// DeleteDocument 删除对应文档的所有chunks
//
// 根据document_id过滤并删除所有相关的向量点
func (qs *QdrantStore) DeleteDocument(documentID string) error {
	ctx := context.Background()

	// 使用fliter删除满足条件的点
	// payload.document_id == documentID
	_, err := qs.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: qs.collection,
		Points: qdrant.NewPointsSelectorFilter(
			&qdrant.Filter{
				Must: []*qdrant.Condition{
					qdrant.NewMatch("document_id", documentID), // 匹配在document_id字段有documentID的值的点
				},
			}),
	})
	if err != nil {
		return fmt.Errorf("删除文档向量失败: %w", err)
	}
	log.Printf("成功删除文档%s的所有向量", documentID)
	return nil
}
