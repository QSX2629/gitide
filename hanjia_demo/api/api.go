package api

import (
	"hanjia_demo/service"
	"hanjia_demo/utils_jwt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Register 注册接口
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}

// Login 登录接口
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, _ := utils_jwt.GenerateToken(user.ID, user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// CreateArticleRequest CreatePost 创建文章/问题
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article, err := service.CreateArticle(userID.(uint), req.Title, req.Content, req.Status, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": article})
}

type FollowRequest struct {
	FollowedID uint `json:"followed_id" binding:"required"` // 要关注的用户ID
}

// FollowUser 关注你的用户
func FollowUser(c *gin.Context) {
	// 1. 获取当前登录用户ID
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 2. 绑定JSON请求体
	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	// 3. 调用服务
	err := service.FollowUser(followerID.(uint), req.FollowedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "关注成功"})
}

// UnfollowUser 取消关注
// id 是要取消关注的用户ID
func UnfollowUser(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 从路径参数获取被取消关注的用户ID
	followedIDStr := c.Param("id")
	followedID, err := strconv.ParseUint(followedIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	err = service.UnfollowUser(followerID.(uint), uint(followedID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "取消关注成功"})
}

// GetFollowingList 获取我的关注列表
func GetFollowingList(c *gin.Context) {
	// 1. 获取当前登录用户ID（鉴权）
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 2. 获取目标用户ID（默认查自己，支持查他人）
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil || targetUserIDStr == "" {
		// 默认查当前登录用户
		targetUserID = uint64(c.GetUint("userID"))
	}

	// 3. 调用服务
	followingList, err := service.GetFollowingList(uint(targetUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": len(followingList),
		"list":  followingList,
	})
}

// GetFollowerList 获取我的粉丝列表
func GetFollowerList(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil || targetUserIDStr == "" {
		targetUserID = uint64(c.GetUint("userID"))
	}

	followerList, err := service.GetFollowerList(uint(targetUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": len(followerList),
		"list":  followerList,
	})
}

type SearchArticleRequest struct {
	Keyword string `form:"keyword" json:"keyword" binding:"required"` // 搜索关键词
	Page    int    `form:"page" json:"page" `                         // 页码
	Size    int    `form:"size" json:"size" `                         // 每页条数
}

// SearchArticle 搜索文章（仅标题/作者）
func SearchArticle(c *gin.Context) {
	// 1. 绑定参数（支持URL查询参数/JSON）
	var req SearchArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	// 2. 校验参数
	if req.Keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 || req.Size > 50 {
		req.Size = 10 // 限制最大条数，防止性能问题
	}

	// 3. 调用服务层
	articles, total, err := service.SearchArticleByTitleOrAuthor(req.Keyword, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败：" + err.Error()})
		return
	}

	// 4. 返回结果（仅包含文章+作者核心信息）
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索成功",
		"data": gin.H{
			"total":    total,    // 总条数
			"page":     req.Page, // 当前页
			"size":     req.Size, // 每页条数
			"articles": articles, // 文章列表（含作者信息）
		},
	})
}

// UpdateArticleRequest 更改文章状态（编辑，发布，删除）
type UpdateArticleRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

func UpdateArticleStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未登录"})
		return
	}
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.UpdateArticleStatus(req.ArticleID, userID.(uint), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章状态修改成功"})
}

// LockUserRequest 禁言用户
type LockUserRequest struct {
	TargetUserID uint `json:"target_user_id" binding:"required"`
}

func LockUser(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未登录"})
	}
	var req LockUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.LockUser(adminID.(uint), req.TargetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "用户已禁言"})
}
func UnlockUser(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未登录"})
		return
	}
	var req LockUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UnlockUser(adminID.(uint), req.TargetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "已取消"})
}

type CommentRequest struct {
	ArticleID uint   `json:"article_id" `
	Content   string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	comment, err := service.CreateComment(userID.(uint), req.ArticleID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment)
}
