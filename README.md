# GO-RAG
尝试使用Go语言编写一个RAG问答系统

## 项目架构
>使用Gin框架搭建后端，在后端中集成RAG功能。前端比较简易，使用Gin框架对前端页面进行静态托管 

技术栈：</br>
- 后端框架：Gin   </br>
- 前端：Gin静态托管HTML文件 </br>
- Docker容器化

RAG：
- 数据分块：手搓分块器
- 向量化Embedding模型：`待定`
- 向量存储：QDurant向量数据库
- Chat模型：`待定`

> 为了实现RAG的纯粹性，先不引入鉴权


