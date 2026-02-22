package middleware_jwt

import (
	"context"
	"errors"
	"fmt"
	"hanjia_demo/redis"
	"hanjia_demo/utils_jwt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

const (
	ArticleLimitCount = 3
	ArticleLimitTime  = 1 * time.Hour
	CommentLimitCount = 7
	CommentLimitTime  = 1 * time.Hour
)

func GetRedisKey(userID uint, action string) string {
	return fmt.Sprintf("rate_limit:%s:%d", userID, action)
}

//action:article/comment
func CheckRateLimit(userID uint, action string) (bool, error) {
	var LimitCount int
	var LimitTime time.Duration
	switch action {
	case "article":
		LimitCount = ArticleLimitCount
		LimitTime = ArticleLimitTime
	case "comment":
		LimitCount = CommentLimitCount
		LimitTime = CommentLimitTime
	default:
		return false, errors.New("无效类型")
	}
	key := GetRedisKey(userID, action)
	//自增器
	count, err := redis.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		redis.RedisClient.Expire(ctx, key, LimitTime)
	}
	if count > int64(LimitCount) {
		return true, nil
	}
	return false, nil
}

// ArticleRateLimit 文章发布防刷中间件
func ArticleRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(200, gin.H{
				"error": "未登录"})
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(200, gin.H{
				"error": "ID格式错误"})
			c.Abort()
			return
		}
		exceed, err := CheckRateLimit(uid, "article")
		if err != nil {
			c.JSON(200, gin.H{
				"error": "检验失败"})
			c.Abort()
			return
		}
		if exceed {
			c.JSON(200, gin.H{
				"error": "频率过高"})
			c.Abort()
			return
		}
		c.Next()

	}
}
func CommentRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(200, gin.H{
				"error": "未登录"})
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(200, gin.H{
				"error": "格式错误"})
			c.Abort()
			return
		}
		exceed, err := CheckRateLimit(uid, "comment")
		if err != nil {
			c.JSON(200, gin.H{
				"error": "限流失败"})
			c.Abort()
			return
		}
		if exceed {
			c.JSON(200, gin.H{
				"error": "超出频率"})
		}
		c.Next()

	}
}
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求头中的token（支持 Bearer 前缀）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		// 处理 Bearer Token 格式
		tokenStr := ""
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenStr = parts[1]
		} else {
			tokenStr = authHeader
		}

		claims, err := utils_jwt.ParseToken(tokenStr)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效或已过期"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
