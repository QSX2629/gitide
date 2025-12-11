package middleware

import (
	"demo/utils"
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
		claims, err := utils.ParseToken(parts[1])
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
