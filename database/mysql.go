package database

import (
	"fmt"
	"go-rag/config"
	"log"
	"time"

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
		return fmt.Errorf("连接SQL数据库失败; %w", err)
	}

	// 获取底层SQL.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败：%w", err)
	}

	// 调整连接池参数
	// 设置空闲连接池中最大的连接数
	sqlDB.SetMaxIdleConns(10)

	// 设置连接最大数量(连接池总连接数上限)
	sqlDB.SetMaxOpenConns(100)

	// 设置连接可复用的最长时间:1h
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 数据表迁移
	err = db.AutoMigrate(&DocumentDB{})
	if err != nil {
		return fmt.Errorf("数据库表迁移失败：%w", err)
	}

	// 打印一下日志
	log.Println("====InitMySQL===========================")
	log.Println("===MySQL数据库连接成功===")
	log.Printf("- 主机： %s:%d", cfg.Host, cfg.Port)
	log.Printf("- 数据库：%s", cfg.Database)
	log.Println("- 表 documents 已经创建")
	log.Println("========================================")

	return nil
}

// GetDB 获取数据库连接实例
//
// 就是获取全局实例db
func GetDB() *gorm.DB {
	return db
}

// CloseMySQL 关闭数据库连接
func CloseMySQL() error {
	// 获取底层sql.DB对象进行close
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 关闭连接
	return sqlDB.Close()

}
