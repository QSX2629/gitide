package service

import (
	"errors"
	"hanjia_demo/dao"
	"hanjia_demo/model"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(username, email, password string) (*model.User, error) {
	// 检查用户是否已存在
	exists, err := dao.CheckUserExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已注册") // 原代码错误：提示"邮箱已注册"，实际检查的是username
	}

	// 密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	newUser := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPwd),
	}
	if err := dao.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login 用户登录
func Login(username, password string) (*model.User, error) {
	// 查询用户
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// CreateArticle 创建文章/问题
func CreateArticle(userID uint, title, content, postType string, status string) (*model.Article, error) {
	// 检查是否被禁言
	adminUser, err := CheckAdminUser(userID)
	if err != nil {
		return nil, err
	}
	if adminUser {
		return nil, errors.New("被禁言")
	}

	// 校验文章状态（原代码的if true是测试代码，这里改为实际校验）
	validStatus := map[string]bool{
		model.ArticleStatusDraft:     true,
		model.ArticleStatusPublished: true,
		model.ArticleStatusDeleted:   true,
	}
	if !validStatus[status] {
		return nil, errors.New("文章状态不合法")
	}

	// 创建文章
	article := &model.Article{
		Title:   title,
		Content: content,
		Status:  status,
		UserId:  userID,
	}
	if err := dao.CreateArticle(article); err != nil {
		return nil, err
	}

	return article, nil
}

// FollowUser 关注用户
func FollowUser(followerID, followedID uint) error {
	// 禁止关注自己
	if followerID == followedID {
		return errors.New("无法关注自己")
	}

	// 创建关注关系
	follow := &model.Follow{
		UserID:     followerID,
		FollowedID: followedID,
	}
	rowsAffected, err := dao.CreateFollow(follow)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("已关注该用户")
	}

	return nil
}

// UnfollowUser 取消关注
func UnfollowUser(followerID, followedID uint) error {
	rowsAffected, err := dao.DeleteFollow(followerID, followedID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("未关注该用户，无需取消")
	}

	return nil
}

// GetFollowingList 获取「我关注的人」列表
func GetFollowingList(userID uint) ([]model.User, error) {
	return dao.GetFollowingList(userID)
}

// GetFollowerList 获取「关注我的人」列表
func GetFollowerList(userID uint) ([]model.User, error) {
	return dao.GetFollowerList(userID)
}

// SearchArticleByTitleOrAuthor 搜索文章
func SearchArticleByTitleOrAuthor(keyword string, page, size int) ([]model.Article, int64, error) {
	return dao.SearchArticleByTitleOrAuthor(keyword, page, size)
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(articleID, userID uint, status string) error {
	// 校验状态合法性
	if status != model.ArticleStatusDeleted && status != model.ArticleStatusDraft && status != model.ArticleStatusPublished {
		return errors.New("文章状态非法")
	}

	// 更新文章状态
	if err := dao.UpdateArticleStatus(articleID, userID, status); err != nil {
		return errors.New("没有文章或没有权限")
	}

	return nil
}

// CheckAdminUser 检查是否是禁言用户（原代码命名有问题，admin_user实际是禁言标识）
func CheckAdminUser(userID uint) (bool, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return false, errors.New("用户不存在")
	}
	return user.Admin_user, nil
}

// CheckCommonUser 检查是否是普通用户
func CheckCommonUser(userID uint) (bool, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user.Common_user, nil
}

// LockUser 禁言用户（原代码逻辑错误：Update的是admin_user=true，实际应该是禁言标识）
func LockUser(adminID, goalID uint) error {
	// 检查操作人是否是管理员
	adminUser, err := CheckAdminUser(adminID)
	if err != nil {
		return err
	}
	if !adminUser {
		return errors.New("无管理员权限")
	}

	// 禁言目标用户（这里假设admin_user=true代表禁言，需根据实际业务调整）
	if err := dao.UpdateUserStatus(goalID, true); err != nil {
		return errors.New("用户不存在")
	}

	return nil
}

// UnlockUser 解禁用户
func UnlockUser(adminID, goalID uint) error {
	// 检查操作人是否是管理员
	adminUser, err := CheckAdminUser(adminID)
	if err != nil {
		return err
	}
	if !adminUser {
		return errors.New("无权限")
	}

	// 解禁目标用户
	if err := dao.UpdateUserStatus(goalID, false); err != nil {
		return errors.New("用户不存在")
	}

	return nil
}

// CreateComment 创建评论
func CreateComment(userID, articleID uint, content string) (*model.Comment, error) {
	// 校验是否被禁言
	adminUser, err := CheckAdminUser(userID)
	if err != nil {
		return nil, err
	}
	if adminUser {
		return nil, errors.New("你已被永久禁言，禁止发布评论")
	}

	// 创建评论
	comment := &model.Comment{
		Content:   content,
		ArticleId: articleID,
	}
	if err := dao.CreateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}
