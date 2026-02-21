package main

import (
	"log"

	"lanshan07/api" // 替换为你的实际模块名
	"lanshan07/db"  // 替换为你的实际模块名
)

func main() {
	// 1. 初始化数据库
	db.Init()

	// 2. 初始化路由
	r := api.InitrouterGin()

	// 3. 启动服务
	log.Println("服务启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
