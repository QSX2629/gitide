package router

import (
	"demo/api"
	"demo/middleware_JWT"
	"demo/middleware_lock"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default() // 带默认中间件（日志+恢复）
	r.Use(middleware_lock.RateLimit(100, 10))

	// 1. 注册/登录/退出路由组（/auth前缀）
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", api.RegisterHandler)
		authGroup.POST("/login", api.LoginHandler)
		authGroup.POST("/logout", api.LogoutHandler)
	}

	// 2. 用户增删改查路由组
	memberGroup := r.Group("/member")
	memberGroup.Use(middleware.JWTAuth())
	memberGroup.Use(middleware.AdminAuth())
	memberGroup.GET("/list", middleware_lock.RateLimit(50, 5), api.ListMembersHandler)
	{
		memberGroup.POST("/add", api.AddMemberHandler)             // 新增
		memberGroup.DELETE("/delete/:id", api.DeleteMemberHandler) // 删除
		memberGroup.PUT("/update/:id", api.UpdateMemberHandler)    // 更新
		memberGroup.GET("/detail/:id", api.GetMemberHandler)       // 详情
	}

	return r
}
