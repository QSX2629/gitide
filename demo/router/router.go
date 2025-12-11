package router

import (
	"demo/api"
	"demo/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default() // 带默认中间件（日志+恢复）

	// 1. 注册/登录路由组（/auth前缀）
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", api.RegisterHandler)
		authGroup.POST("/login", api.LoginHandler)
	}

	// 2. 会员增删改查路由组（/member前缀）
	memberGroup := r.Group("/member")
	memberGroup.Use(middleware.JWTAuth())
	memberGroup.Use(middleware.AdminAuth())
	{
		memberGroup.POST("/add", api.AddMemberHandler)             // 新增
		memberGroup.DELETE("/delete/:id", api.DeleteMemberHandler) // 删除
		memberGroup.PUT("/update/:id", api.UpdateMemberHandler)    // 更新
		memberGroup.GET("/list", api.ListMembersHandler)           // 列表
		memberGroup.GET("/detail/:id", api.GetMemberHandler)       // 详情
	}

	return r
}
