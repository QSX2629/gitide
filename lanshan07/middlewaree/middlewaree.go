package middlewaree

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
