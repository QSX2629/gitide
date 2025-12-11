package api

import (
	"lanshan07/lanshan07/db"
	"lanshan07/lanshan07/middlewaree"
	"lanshan07/lanshan07/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitrouterGin 路由编写
func InitrouterGin() *gin.Engine {
	r := gin.Default()
	r.Use(middlewaree.CorsDefault())
	memberGroup := r.Group("/api/member")
	{
		memberGroup.POST("/", createMember)      //增路由
		memberGroup.GET("/:id", getMemberID)     //查路由
		memberGroup.PUT("/:id", updateMember)    //改路由
		memberGroup.DELETE("/:id", deleteMember) //删路由
	}
	return r
}

// /////////////////********创建********//////////////////////////
func createMember(c *gin.Context) {
	var member model.Member
	//绑定请求体于结构体
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数失败"})
		return
	}
	//存入数据库
	if err := db.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建成员失败"})
	}
	//返回成功
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    member,
	})
}

// //////////////********查找*************///////////////
func getMemberID(c *gin.Context) {
	//从URL获取id，并且将代转换的id参数转为10进制的64位整数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID格式错误",
		})
		return
	}
	var member model.Member
	//从数据库中查询
	if err := db.DB.First(&member, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "会员不存在",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "查询失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"data":    member,
	})
	return
}

// ************更新****************///////////////////
func updateMember(c *gin.Context) {
	idStr := c.Param("id") //获取id
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID格式错误",
		})
		return
	}
	//绑定更新参数
	var updateData model.Member
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 先查询会员是否存在
	var member model.Member
	if err := db.DB.First(&member, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "会员不存在"})
			return
		}
		c.JSON(500, gin.H{"error": "查询会员失败: " + err.Error()})
		return
	}
	//更新执行
	if err := db.DB.Model(&member).Updates(updateData).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "更新失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// deleteMember 删除会员
func deleteMember(c *gin.Context) {
	// 获取ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID格式错误"})
		return
	}

	// 执行删除
	if err := db.DB.Delete(&model.Member{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "删除会员失败: " + err.Error()})
		return
	}

	// 返回删除成功
	c.JSON(200, gin.H{"message": "会员删除成功"})
}
