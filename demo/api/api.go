package api

import (
	"context"
	"demo/service"
	"demo/utils_JWT"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required,min=6,max=10"`
		Major    string `json:"major" binding:"required"`
	}
	// 参数校验
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误：" + err.Error()})
		return
	}
	// 调用业务逻辑
	if err := service.Register(req.Account, req.Password, req.Major); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()}) // 直接返回错误字符串
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	member, err := service.Login(req.Account, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	expirationTime := time.Now().Add(2 * time.Hour)
	token, err := utils.GenerateToken(member.Account, member.Character, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "令牌生成失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{"id": strconv.FormatUint(uint64(member.ID), 10),
			"account":   member.Account,
			"token":     token,
			"character": member.Character,
		},
	})
}
func AddMemberHandler(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required,min=6,max=10"`
		Major    string `json:"major" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	if err := service.AddMember(req.Account, req.Password, req.Major); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "新增会员成功"})
}

func DeleteMemberHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID格式错误"})
		return
	}
	if err := service.DeleteMember(uint(id)); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除用户成功"})
}

func UpdateMemberHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID格式错误"})
		return
	}
	var req struct {
		NewPassword string `json:"new_password"`
		Major       string `json:"major" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	if err := service.UpdateMember(uint(id), req.NewPassword, req.Major); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新用户成功"})
}

func ListMembersHandler(c *gin.Context) {
	members, err := service.ListMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败：" + err.Error()})
		return
	}
	var respData []gin.H
	for _, m := range members {
		respData = append(respData, gin.H{
			"id":         strconv.FormatUint(uint64(m.ID), 10),
			"account":    m.Account,
			"major":      m.Major,
			"created_at": m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": respData})
}

func GetMemberHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID格式错误"})
		return
	}
	member, err := service.GetMemberByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":         strconv.FormatUint(uint64(member.ID), 10),
			"account":    member.Account,
			"major":      member.Major,
			"created_at": member.CreatedAt.Format("2006-01-02 15:04:05"), // 时间格式化（可选优化）
		},
	})
}
func LogoutHandler(c *gin.Context) {
	//获取token
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "没有令牌",
		})
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "") {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "令牌格式错误",
		})
		return
	}
	tokenStr := parts[1]
	//解析令牌
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "令牌无效",
		})
		return
	}
	//加入黑名单
	ctx := context.Background()
	expireTime := time.Until(claims.ExpiresAt.Time)
	_ = utils.AddTokenToBlackList(ctx, tokenStr, expireTime)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"meesage": "推出登录",
	})
}
