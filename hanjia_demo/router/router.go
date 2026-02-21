package router

import (
	"hanjia_demo/api"
	"hanjia_demo/middleware_jwt"
	"hanjia_demo/middware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 公开接口
	public := r.Group("/api")
	{
		public.POST("/register", api.Register)
		public.POST("/login", api.Login)
		public.GET("/articles/search", api.SearchArticle)
	}

	// 私有接口（需要登录）
	private := r.Group("/api")
	private.Use(middleware_jwt.JWTAuth())
	{
		private.POST("/article", middleware_jwt.ArticleRateLimit(), api.CreateArticle)
		private.POST("/comment", middleware_jwt.CommentRateLimit(), api.CreateComment)
	}
	private.Use(middware.AuthRequired())
	{
		private.POST("/admin/lock", api.LockUser)
		private.POST("/admin/unlock", api.UnlockUser)
		private.POST("/posts", api.CreateArticle)
		private.POST("/follow", api.FollowUser)
		private.DELETE("/follow/:id", api.UnfollowUser)
		private.GET("/users/following", api.GetFollowingList)
		private.GET("/users/followers", api.GetFollowerList)
		private.POST("/articles/status", api.UpdateArticleStatus)

	}

	return r
}
