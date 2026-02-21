package main

import (
	"hanjia_demo/db"
	"hanjia_demo/redis"
	"hanjia_demo/router"
	"hanjia_demo/utils_log"
	"hanjia_demo/utils_viper"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	// 1. 先初始化配置
	if err := utils_viper.InitConfig(); err != nil {
		panic("配置初始化失败: " + err.Error())
	}

	// 2. 初始化日志
	if err := utils_log.InitLogger(); err != nil {
		panic("日志初始化失败: " + err.Error())
	}
	defer utils_log.Logger.Sync()

	// 3. 初始化MySQL
	db.Init()

	// 4. 初始化Redis
	redis.InitRedis()

	// 5. 初始化路由
	r := router.SetupRouter()

	// 6. 启动服务
	port := utils_viper.GetConfig().App.Port
	utils_log.Info("服务器启动成功", zap.Int("port", port))
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		utils_log.Fatal("服务器启动失败", zap.Error(err))
	}
}
