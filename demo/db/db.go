package db

import (
	"demo/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局DB实例，供其他包调用

// InitDB 初始化数据库连接
func InitDB() error {
	// 替换为你的MySQL配置：用户名:密码@tcp(地址:端口)/数据库名?参数
	dsn := "root:Qsx20061212@tcp(127.0.0.1:3306)/moon?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err // 连接失败返回错误
	}
	// 自动创建members表（不存在则创建，存在则同步结构）
	return DB.AutoMigrate(&model.Member{})
}
