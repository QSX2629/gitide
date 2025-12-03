//CorsDefault（跨域请求处理）和 JWT（身份验证）。中间件的本质是「请求到达接口前的 “拦截器”」
//，用于统一处理通用逻辑（如跨域、权限验证），避免在每个接口重复写相同代码。下面逐段拆解，讲清作用、原理和细节：

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
		// 修正后：先解析出 Claims 结构体，再取 Username 字段
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "令牌无效或已过期：" + err.Error()})
			c.Abort()
			return
		}
		// 将用户名存入上下文（取 claims.Username）
		c.Set("username", claims.Username)
		// 继续执行后续接口逻辑
		c.Next()
	}
}
