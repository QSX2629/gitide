package middleware

import (
	"context"
	"demo/rediss"
	"demo/utils_JWT"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "没有令牌",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "令牌格式错误",
			})
			c.Abort()
			return
		}
		tokenStr := parts[1]
		//检查令牌是否在黑名单中
		ctx := context.Background()
		blacklistKey := fmt.Sprintf("blacklist:%s", tokenStr)
		exists := rediss.ExistsCache(ctx, blacklistKey)
		if exists {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "令牌以注销",
			})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("account", claims.Account)
		c.Set("character", claims.Character)
		c.Next()
	}
}
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		character, exists := c.Get("character")
		if !exists || character != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "没有权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
