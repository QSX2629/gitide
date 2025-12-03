package api

import (
	"lanshan06/dao"
	"lanshan06/middleware"
	"lanshan06/model"
	"lanshan06/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
	}
	if dao.CheckUserExists(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user already exists",
		})
		return
	}
	if err := dao.AddUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "register failed",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}
	if !dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found",
		})
		return
	}
	token, err := utils.GenerateToken(req.Username, time.Now().Add(10*time.Minute))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "login",
	})
}

// Ping1 假设在 api.go 中定义（需与路由注册处同包，或通过包导入）
func Ping1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func InitrouterGin() {
	r := gin.Default()

	// 全局注册跨域中间件（所有接口支持跨域）
	r.Use(middleware.CorsDefault())

	// 1. 公开接口（无需登录）
	publicGroup := r.Group("/api/public")
	{
		publicGroup.POST("/login", Login)       // 登录
		publicGroup.POST("/register", Register) // 注册
	}

	// 2. 需登录接口（JWT 中间件保护）
	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.JWT()) // 该组下所有接口需验证令牌
	{
		authGroup.GET("/ping", Ping1)                      // 原有需要登录的接口
		authGroup.POST("/modify-password", ModifyPassword) // 新增：修改密码接口
	}

	// 注意：r.Run() 必须在所有路由注册之后！
	r.Run(":8080")

}
func ModifyPassword(c *gin.Context) {
	// 1. 绑定请求参数（需在 model 中定义 ModifyPasswordRequest 结构体）
	var req model.ModifyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误：" + err.Error()})
		return
	}

	// 2. 从 JWT 上下文获取登录用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取用户信息失败"})
		return
	}
	userNameStr := username.(string)

	// 3. 从 dao 层获取加密后的旧密码
	oldHashedPwd := dao.SelectPasswordFromUsername(userNameStr)
	if oldHashedPwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户不存在"})
		return
	}

	// 4. 校验旧密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(oldHashedPwd), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "旧密码输入错误"})
		return
	}

	// 5. 调用 dao 层更新密码（dao 层自动加密新密码）
	if err := dao.UpdatePassword(userNameStr, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "修改密码失败"})
		return
	}

	// 6. 响应成功
	c.JSON(http.StatusOK, gin.H{"message": "修改密码成功，请重新登录"})
}
