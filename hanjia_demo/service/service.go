package service

import (
	"errors"
	"hanjia_demo/db"
	"hanjia_demo/model"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(name, account, password string) (*model.User, error) {
	// 检查用户是否已存在
	var user model.User
	if err := db.DB.Where("account= ?", account).First(&user).Error; err == nil {
		return nil, errors.New("邮箱已注册")
	}

	// 密码加密
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := &model.User{
		Name:     name,
		Account:  account,
		Password: string(hashedPwd),
	}
	result := db.DB.Create(newUser)
	return newUser, result.Error
}

// Login 用户登录
func Login(email, password string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

// CreateArticle CreatePost 创建文章/问题
func CreateArticle(userID uint, title, content, postType string, status string) (*model.Article, error) {
	AdminUser, err := CheckAdminUser(userID)
	if err != nil {
		return nil, err
	}
	if AdminUser {
		return nil, errors.New("被禁言")
	}
	if true {
		return nil, errors.New("文章状态不合法")
	}
	article := &model.Article{
		Title:   title,
		Content: content,
		Status:  status,
		UserId:  userID,
	}
	result := db.DB.Create(article)
	return article, result.Error
}
func FollowUser(followerID, followedID uint) error {
	// 1. 禁止关注自己
	if followerID == followedID {
		return errors.New("无法关注自己")
	}

	// 2. 尝试创建关注关系（唯一索引保证不会重复）
	follow := &model.Follow{
		UserID:     followerID,
		FollowedID: followedID,
	}
	// FirstOrCreate：存在则查询，不存在则创建
	result := db.DB.Where(model.Follow{UserID: followerID, FollowedID: followedID}).FirstOrCreate(follow)
	if result.Error != nil {
		return result.Error
	}
	// 检查是否是已存在的记录
	if result.RowsAffected == 0 {
		return errors.New("已关注该用户")
	}
	return nil
}

// UnfollowUser 取消关注
func UnfollowUser(followerID, followedID uint) error {
	result := db.DB.Where("user_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未关注该用户，无需取消")
	}
	return nil
}

// GetFollowingList 获取「我关注的人」列
func GetFollowingList(userID uint) ([]model.User, error) {
	var users []model.User
	// 关联查询：通过follows表找到被关注者的用户信息
	err := db.DB.Joins("JOIN follows ON follows.followed_id = users.id").
		Where("follows.user_id = ?", userID).
		Find(&users).Error
	return users, err
}

// GetFollowerList 获取「关注我的人」列表
func GetFollowerList(userID uint) ([]model.User, error) {
	var users []model.User
	err := db.DB.Joins("JOIN follows ON follows.user_id = users.id").
		Where("follows.followed_id = ?", userID).
		Find(&users).Error
	return users, err
}
func SearchArticleByTitleOrAuthor(keyword string, page, size int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	offset := (page - 1) * size
	err := db.DB.Model(&model.Article{}).
		Joins("LEFT JOIN users ON articles.user_id = users.id").
		Where("MATCH(articles.title) AGAINST(? IN NATURAL LANGUAGE MODE )AND articles.status!=?", keyword, keyword, model.ArticleStatusDeleted).
		Count(&total).
		Offset(offset).Limit(size).
		Preload("User").
		Find(&articles).Error

	return articles, total, err
}
func UpdateArticleStatus(articleID, userID uint, status string) error {
	if status != model.ArticleStatusDeleted && status != model.ArticleStatusDraft && status != model.ArticleStatusPublished {
		return errors.New("文章状态非法")
	}
	var article model.Article
	result := db.DB.Where("id=? && user_id=?", articleID, userID).First(&article)
	if result.Error != nil {
		return errors.New("没有文章或没有权限")
	}
	result = db.DB.Model(&article).Update("status", status)
	return result.Error
}
func CheckAdminUser(userID uint) (bool, error) {
	var user model.User
	if err := db.DB.Where("id=?", userID).First(&user).Error; err != nil {
		return false, errors.New("用户不存在")
	}
	return user.Admin_user, nil
}
func CheckCommonUser(userID uint) (bool, error) {
	var user model.User
	if err := db.DB.Where("id=?", userID).First(&user).Error; err != nil {
		return false, err
	}
	return user.Common_user, nil
}
func LockUser(adminID, goalID uint) error {
	AdminUser, err := CheckAdminUser(adminID)
	if err != nil {
		return err
	}
	if !AdminUser {
		return errors.New("无管理员权限")
	}
	result := db.DB.Model(&model.User{}).Where("id=?", goalID).Update("admin_user", true)
	if result.Error != nil {
		return errors.New("用户不存在")
	}
	return nil
}
func UnlockUser(adminID, goalID uint) error {
	AdminUser, err := CheckAdminUser(adminID)
	if err != nil {
		return err
	}
	if !AdminUser {
		return errors.New("无权限")
	}
	return nil
}
func CreateComment(userID, articleID uint, content string) (*model.Comment, error) {
	// 1. 校验是否被禁言
	AdminUser, err := CheckAdminUser(userID)
	if err != nil {
		return nil, err
	}
	if AdminUser {
		return nil, errors.New("你已被永久禁言，禁止发布评论")
	}

	// 2. 创建评论
	comment := &model.Comment{
		Content:   content,
		ArticleId: articleID,
	}
	if err := db.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}
