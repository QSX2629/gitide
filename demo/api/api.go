package api

import (
	"demo/service"
	"demo/utils"
	"net/http"
	"strconv"
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
		NewPassword string `json:"new_password"` // 可选更新
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
