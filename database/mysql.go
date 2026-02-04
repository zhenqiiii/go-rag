package database

import (
	"fmt"
	"go-rag/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库连接实例
var db *gorm.DB

// InitMySQL 初始化MySQL数据库连接
//
// 传入MySQLConfig进行配置，返回error
func InitMySQL(cfg config.MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败; %w", err)
	}
	return nil
}
