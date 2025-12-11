package main

import (
	"demo/db"
	"demo/router"
	"demo/service"
)

func main() {
	// 1. 初始化数据库连接
	if err := db.InitDB(); err != nil {
		panic("数据库初始化失败：" + err.Error())
	}
	if err := service.CreateDefaultAdmin(); err != nil {
		panic("默认管理员创建失败：" + err.Error())
	}
	// 打印提示
	println("默认管理员创建成功（账号：admin，密码：Admin123!）")
	// 2. 配置路由
	r := router.SetupRouter()
	// 3. 启动服务（监听8080端口）
	if err := r.Run(":8080"); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}
