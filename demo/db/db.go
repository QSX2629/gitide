package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局DB实例，供其他包调用

// InitDB 初始化数据库连接
func InitDB() error {
	dsn := "root:Qsx20061212@tcp(127.0.0.1:3306)/moon?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err // 连接失败返回错误
	}
	return nil
}
