package dao

import (
	"hanjia_demo/db"
	"hanjia_demo/model"
)

// ====================== 用户相关数据库操作 ======================
// CheckUserExists 检查用户名是否存在
func CheckUserExists(username string) (bool, error) {
	var user model.User
	err := db.DB.Where("username= ?", username).First(&user).Error
	if err != nil {
		// 如果是记录不存在的错误，返回false；否则返回错误
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateUser 创建新用户
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID查询用户
func GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := db.DB.Where("id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus 更新用户状态（禁言/解禁）
func UpdateUserStatus(userID uint, adminUser bool) error {
	result := db.DB.Model(&model.User{}).Where("id=?", userID).Update("admin_user", adminUser)
	return result.Error
}

// ====================== 文章相关数据库操作 ======================
// CreateArticle 创建文章
func CreateArticle(article *model.Article) error {
	return db.DB.Create(article).Error
}

// SearchArticleByTitleOrAuthor 按标题/作者搜索文章（带分页）
func SearchArticleByTitleOrAuthor(keyword string, page, size int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * size

	// 构建查询
	query := db.DB.Model(&model.Article{}).
		Joins("LEFT JOIN users ON articles.user_id = users.id").
		Where("MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE )AND status!=?", keyword, model.ArticleStatusDeleted)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询并预加载用户信息
	err := query.
		Offset(offset).Limit(size).
		Preload("User").
		Find(&articles).Error

	return articles, total, err
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(articleID, userID uint, status string) error {
	var article model.Article
	// 先查询文章是否存在且属于当前用户
	result := db.DB.Where("id=? AND user_id=?", articleID, userID).First(&article)
	if result.Error != nil {
		return result.Error
	}
	// 更新状态
	return db.DB.Model(&article).Update("status", status).Error
}

// GetArticleByID 根据ID查询文章
func GetArticleByID(articleID uint) (*model.Article, error) {
	var article model.Article
	err := db.DB.Where("id=?", articleID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateFollow 创建关注关系
func CreateFollow(follow *model.Follow) (int64, error) {
	result := db.DB.Where(model.Follow{UserID: follow.UserID, FollowedID: follow.FollowedID}).FirstOrCreate(follow)
	return result.RowsAffected, result.Error
}

// DeleteFollow 取消关注
func DeleteFollow(followerID, followedID uint) (int64, error) {
	result := db.DB.Where("user_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Follow{})
	return result.RowsAffected, result.Error
}

// GetFollowingList 获取关注列表
func GetFollowingList(userID uint) ([]model.User, error) {
	var users []model.User
	err := db.DB.Joins("JOIN follows ON follows.followed_id = users.id").
		Where("follows.user_id = ?", userID).
		Find(&users).Error
	return users, err
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(userID uint) ([]model.User, error) {
	var users []model.User
	err := db.DB.Joins("JOIN follows ON follows.user_id = users.id").
		Where("follows.followed_id = ?", userID).
		Find(&users).Error
	return users, err
}

// CreateComment 创建评
func CreateComment(comment *model.Comment) error {
	return db.DB.Create(comment).Error
}
