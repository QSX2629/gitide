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

// Register 用户注册接口  核心流程：接收注册请求  +校验参数 + 检查用户名是否已存在+ 新增用户到数据库 + 返回结果
func Register(c *gin.Context) {
	// 1. 绑定前端JSON请求参数到 model.User 结构体（参数校验）
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// 2. 调用 dao 层检查用户名是否已存在（数据库查询）
	if dao.CheckUserExists(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "service already exists"})
		return
	}

	// 3. 调用 dao 层新增用户（数据库插入，dao层内部会对密码加密）
	if err := dao.AddUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "register failed"})
		return
	}

	// 4. 注册成功
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// Login 接收登录请求 → 校验参数 → 验证用户名 + 密码是否匹配 → 生成 JWT 令牌 → 返回令牌给前端
func Login(c *gin.Context) {
	// 1. 绑定前端JSON参数（用户名+密码）
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// 2. 调用 dao 层验证用户名和密码（dao层会拿明文密码和数据库中的哈希密码比对）
	if !dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "service not found"})
		return
	}

	// 3. 生成JWT令牌（有效期10分钟）
	token, err := utils.GenerateToken(req.Username, time.Now().Add(30*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// 4. 登录成功，返回令牌（前端后续请求需携带该令牌）
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "login"})
}

// Ping1 登录后测试接口:验证 JWT 认证是否生效的测试接口，仅登录后（携带有效令牌）可访问。
func Ping1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// ModifyPassword 修改密码接口:登录后 → 接收旧密码 + 新密码 → 校验旧密码正确性 → 更新数据库中的密码 → 返回结果
func ModifyPassword(c *gin.Context) {
	//绑定请求参数（旧+新，在model中的ModifyPasswordRequest）
	var req model.ModifyPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 核心修改：返回具体错误信息，而非模糊的"bad request"
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"error":   err.Error(), // 暴露真实错误
		})
		return
	}
	//JWT上下文获取当前的用户名（通过中间件后存入  ）
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "service not found"})
		return
	}
	usernamestring := username.(string)
	//从数据库得到旧的哈希密码
	oldHashedPassword := dao.SelectPasswordFromUsername(usernamestring)
	if oldHashedPassword == " " {
		c.JSON(http.StatusBadRequest, gin.H{"message": "service can not found"})
		return
	}
	//检验密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(oldHashedPassword), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "old password error"})
		return
	}
	//调用dao层更改密码
	if err := dao.UpdatePassword(usernamestring, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "update password error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
func InitrouterGin() {
	r := gin.Default()
	r.Use(middleware.CorsDefault())
	publicGroup := r.Group("/api/public")
	{
		publicGroup.POST("/login", Login)
		publicGroup.POST("/ping", Register)
		authGroup := r.Group("/api/auth")
		authGroup.Use(middleware.JWT())
		{
			authGroup.GET("/ping", Ping1)
			authGroup.POST("/modify-password", ModifyPassword)
		}
		r.Run(":8080")
	}
}
