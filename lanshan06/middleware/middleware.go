package middleware

import (
	"lanshan06/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 固定跨域配置（开发环境适用）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// JWT 简化版JWT中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取令牌（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未提供令牌"})
			c.Abort() // 终止请求
			return
		}

		// 2. 提取token字符串（去掉"Bearer "前缀）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "令牌格式错误（应为Bearer <token>）"})
			c.Abort()
			return
		}
		tokenStr := parts[1]

		// 3. 解析令牌
		username, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "令牌无效或已过期：" + err.Error()})
			c.Abort()
			return
		}

		// 4. 令牌有效，将用户名存入上下文（后续接口可直接获取）
		c.Set("username", username)

		// 继续执行后续接口逻辑
		c.Next()
	}
}
