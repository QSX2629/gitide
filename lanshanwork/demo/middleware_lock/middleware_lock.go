package middleware_lock

import (
	"demo/utils_lock"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimit 接口限流中间件
func RateLimit(capacity int, rate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 1. 定义限流标识
		var limitKey string
		// 从上下文获取登录用户的账号
		account, exists := c.Get("account")
		if exists {
			limitKey = fmt.Sprintf("account:%s:%s", account, c.FullPath())
		} else {
			// 获取客户端IP
			ip := c.ClientIP()
			limitKey = fmt.Sprintf("ip:%s:%s", ip, c.FullPath())
		}

		// 2. 执行限流判断
		allow, _, err := utils_lock.RateLimit(ctx, limitKey, capacity, rate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "限流服务异常"})
			c.Abort()
			return
		}

		// 3. 拒绝请求（限流触发）
		if !allow {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "msg": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}

		// 4. 允许请求，继续执行
		c.Next()
	}
}
