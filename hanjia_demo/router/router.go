package router

import (
	"hanjia_demo/api"
	"hanjia_demo/middleware_jwt"
	"hanjia_demo/middware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并返回 Gin 引擎
func SetupRouter() *gin.Engine {
	// 使用默认引擎（包含 Logger 和 Recovery 中间件）
	r := gin.Default()

	// 根 API 分组
	apiGroup := r.Group("/api")
	{
		//  公开接口（无需登录）
		publicGroup := apiGroup.Group("/public")
		{
			publicGroup.POST("/register", api.Register)             // 注册
			publicGroup.POST("/login", api.Login)                   // 登录
			publicGroup.POST("/articles/search", api.SearchArticle) // 文章搜索
		}

		//  私有接口（需要登录验证）
		privateGroup := apiGroup.Group("/private")
		// 先挂载 JWT 认证中间件（所有私有接口都需要登录）
		privateGroup.Use(middleware_jwt.JWTAuth())
		{
			// 基础业务接口（仅需登录）
			privateGroup.POST("/articles", middleware_jwt.ArticleRateLimit(), api.CreateArticle) // 创建文章
			privateGroup.POST("/comments", middleware_jwt.CommentRateLimit(), api.CreateComment) // 创建评论
			privateGroup.POST("/follow", api.FollowUser)                                         // 关注用户
			privateGroup.DELETE("/follow/:id", api.UnfollowUser)                                 // 取消关注
			privateGroup.GET("/users/following", api.GetFollowingList)                           // 获取关注列表
			privateGroup.GET("/users/followers", api.GetFollowerList)                            // 获取粉丝列表
			privateGroup.POST("/articles/status", api.UpdateArticleStatus)                       // 更新文章状态

			// 管理员接口（需要额外的权限校验）
			adminGroup := privateGroup.Group("/admin")
			adminGroup.Use(middware.AuthRequired()) // 管理员权限校验
			{
				adminGroup.POST("/lock", api.LockUser)     // 锁定用户
				adminGroup.POST("/unlock", api.UnlockUser) // 解锁用户
			}
		}
	}

	return r
}
