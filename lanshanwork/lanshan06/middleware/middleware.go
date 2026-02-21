package middleware

import (
	"lanshan06/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("【原始Authorization】：%q", authHeader) // %q 会显示空格/空字符，更直观
		log.Printf("【头长度】：%d", len(authHeader))

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "没有令牌"})
			c.Abort()
			return
		}

		// 2. 拆分并打印拆分结果
		parts := strings.SplitN(authHeader, " ", 2) // SplitN 最多拆2段，避免多空格问题
		log.Printf("【拆分结果】：%v，长度：%d", parts, len(parts))
		log.Printf("【拆分后第一段】：%q", parts[0]) // 看是否是 "Bearer"

		// 3. 校验逻辑
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "格式错误"})
			c.Abort()
			return
		}

		//解析令牌
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "令牌过期" + err.Error(),
			})
			c.Abort()
			return
		}
		//令牌有效，存储到gin上下文
		c.Set("username", claims.Username)
		c.Next()
	}
}
func CorsDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		//设置跨域允许的响应头
		c.Header("Access-Control-Allow-Origin", "*")                                         //允许所有的跨域名跨域
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")              //HTTP方法
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,Accept") //请求头
		c.Header("Access-Control-Allow-Credentials", "true")                                 //允许携带cookie
		//处理请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
