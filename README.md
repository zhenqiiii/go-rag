# GO-RAG
尝试使用Go语言编写一个RAG问答系统

## 项目架构
>使用Gin框架搭建后端，在后端中集成RAG功能。前端比较简易，使用Gin框架对前端页面进行静态托管 

技术栈：</br>
- 后端框架：Gin   </br>
- 前端：Gin静态托管HTML文件 </br>
- Docker容器化

RAG：
- 数据分块：`gojieba`
- 向量化Embedding模型：`待定`
- 向量存储：Qdrant向量数据库
- Chat模型：`待定`

> 为了实现的RAG的纯粹性，先不引入鉴权


## 实现
### RAG
#### chunker
分割策略:固定长度切分,同时在标点符号处切分 

测试效果: <br/>
![chunker效果测试](asset/chunk_result.png)

#### embedder

分为`embeddingClient`和`embedder`两个部分

`embeddingClient`：负责发送embedding请求，接收响应并返回 <br/>

`embedder`：是整个pipeline中的embed组件，负责接收前一节点的文本，调用`embeddingClient`进行向量化，并将向量结果发送给下一节点 <br/>

测试效果：
![调用BGE-m3返回1024维向量](asset/embed_test.png)


#### store
接入Qdrant向量库实现 <br/>

目前的搜索逻辑非常简单：基于余弦相似度，返回`top-K`个分数最高的向量
```
文档入库流程：
Document → Chunker → Chunks → Embedder → Vectors → Qdrant.AddChunks()
                                                         ↓
                                                    Points (id, vector, payload)

查询流程：
Query → Embedder → QueryVector → Qdrant.Search() → Top-K Results
                                                      ↓
                                         RetrievalResult[chunk, score]

删除流程：
DocumentID → Qdrant.Delete() → Filter by document_id → Delete Points
```

测试效果：
```go
query := "什么是Golang"
	embeddedQuery, err := embedder.EmbedQuery(query)
	if err != nil {
		log.Fatalf("向量化query失败: %v", err)
	}
	results, err := store.Search(embeddedQuery, 3)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
```
![search结果](/asset/retrieve_result.png)


#### generator
调用的是`Qwen3-8B`，但是感觉因为之前进行过两次重复的`Golang.txt`文件的embedding，导致检索结果有重叠，使得参考的内容灰常有限

```go
query := "什么是Golang"
```

![generate测试结果](/asset/generate_result.png)

```go
query := "Golang的应用场景有哪些？"
```

![generate测试结果2](/asset/generate_result2.png)


#### pipeline
把模型换成了`glm-4.7-flashx`,然后简单试了一下pipeline的效果:

![pipeline的最终效果](/asset/pipeline_test.png)

感觉还是有很大问题的,过程中碰到了些小插曲,然后`maxTokens`调了一下才正

常显示这些内容,输出中提到信息中断,那肯定`chunker`组件也有问题,或者是

`scoreThreshold`出问题了?


![调试](/asset/pipeline_test2.png)

这次把整个`RAGResponse`结构体都打印了出来,`usedChunks`是6没问题,输出的

`chunk`也没啥问题,可能是`maxTokens`设太小了?但是我设了2000的,应该完全

够输出啊? 晕

破案了:切分的问题,给的6个`chunk`里面没有完整信息哈哈,多给几个试试

![找到问题](/asset/pipeline_test3.png)

解决方法的话,调大`chunkSize`或者`topK`应该都可以，主要是信息没给全

## 其他

### 关于余弦相似度
通过两个向量夹角的余弦来度量二者的相似性。

对于一个角的余弦值来说，范围在-1 ~ 1 之间

两个向量方向一致，则余弦相似度为1；正交，则余弦相似度为0；方向完全相反，相似度为-1


余弦相似度关注方向，而非幅度，和文本嵌入是非常吻合的。

Qdrant返回的分数也是如此，在-1 ~ 1之间
![qdrant文档给出的说明](/asset/qdrant_cosine.png)

> 一开始看[wiki](https://zh.wikipedia.org/wiki/%E4%BD%99%E5%BC%A6%E7%9B%B8%E4%BC%BC%E6%80%A7)以为是0~1（正向空间）来着，没想到


