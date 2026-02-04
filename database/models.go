package database

import (
	"time"

	"gorm.io/gorm"
)

// DocumentDB 文档的数据库结构体:专门负责数据库这块的操作
//
// 并不打算存储内容Content
//
// (如果后续想写的话可以加,但是有个问题是不同类型的文件如何处理)
type DocumentDB struct {
	ID         string `gorm:"column:id;type:varchar(50);primaryKey" json:"id"`
	Filename   string `gorm:"column:filename;type:varchar(255);not null" json:"filename"`
	ChunkCount int    `gorm:"column:chunk_count;type:int;default:0" json:"chunk_count"`
	// gorm字段
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;autoCreateTime" json:"created_at"`
	// 开启gorm的软删除特性,执行删除操作后设置该字段,后续的方法找不到该记录
	Deleted gorm.DeletedAt
}

// TableName 指定DocumentDB表名
func (DocumentDB) TableName() string {
	return "documents"
}

// <==========数据库模型的CRUD操作=============>

// GetDocumentByID 顾名思义
//
// 接收文档ID，返回DocumentDB对象+error
func GetDocumentByID(id string) (*DocumentDB, error) {
	var document DocumentDB
	// 查询操作
	// 未查询到记录时，First方法返回ErrorRecordNotFound错误
	if err := db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}

	// 查询成功，返回文档对象
	return &document, nil
}

// GetAllDocuments 获取所有文档的列表
//
// 没有入参，返回DocumentDB的slice+error
func GetAllDocuments() ([]DocumentDB, error) {
	var documents []DocumentDB
	// 直接检索全部对象(Find方法没有找到记录时不会报错)
	if err := db.Find(&documents).Error; err != nil {
		return nil, err
	}
	// 检索成功
	return documents, nil
}

// DeleteDocument 删除文档（由gorm执行软删除）
//
// 入参： 文档的ID
// 返回： error
func DeleteDocument(id string) error {
	// // 通过主键删除
	// db.Delete(&DocumentDB{}, id)
	// 或者是这样声明一个匿名对象进行匹配删除
	if err := db.Delete(&DocumentDB{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
