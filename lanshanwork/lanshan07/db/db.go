package db

import (
	"fmt"
	"lanshan07/model"
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局数据库实例，供其他包调用

// Init 初始化数据库连接
func Init() {
	rawPwd := "Qsx20061212" // 替换为你的MySQL密码
	encodedPwd := url.QueryEscape(rawPwd)

	// 数据库连接DSN
	dsn := "root:%s@tcp(127.0.0.1:3306)/moon?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=false&tls=preferred"
	dsn = fmt.Sprintf(dsn, encodedPwd)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移表结构
	err = DB.AutoMigrate(&model.Member{})
	if err != nil {
		log.Fatalf("表结构迁移失败: %v", err)
	}
	log.Println("表结构迁移成功")
}

//////////////*************用于数据库连接******************////////////////
