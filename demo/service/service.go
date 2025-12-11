package service

import (
	"demo/db"
	"demo/model"
	"errors"
)

// Register 用户注册
func Register(account, password, major string) error {
	// 检查账号是否已存在
	var exist model.Member
	if err := db.DB.Where("account = ?", account).First(&exist).Error; err == nil {
		return errors.New("账号已存在") // 直接返回错误字符串
	}
	// 创建用户（BeforeCreate自动加密密码）
	return db.DB.Create(&model.Member{
		Account:  account,
		Password: password,
		Major:    major,
	}).Error
}

// Login 用户登录
func Login(account, password string) (*model.Member, error) {
	var member model.Member
	// 查询账号
	if err := db.DB.Where("account = ?", account).First(&member).Error; err != nil {
		return nil, errors.New("账号不存在")
	}
	// 校验密码
	if !member.CheckPassword(password) {
		return nil, errors.New("密码错误")
	}
	return &member, nil
}

// AddMember 增加用户
func AddMember(account, password, major string) error {
	return Register(account, password, major)
}

// DeleteMember 根据ID删除用户
func DeleteMember(id uint) error {
	var member model.Member
	if err := db.DB.First(&member, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	return db.DB.Delete(&member).Error
}

// UpdateMember 根据ID更新用户信息
func UpdateMember(id uint, newPwd, major string) error {
	var member model.Member
	if err := db.DB.First(&member, id).Error; err != nil {
		return errors.New("用户不存在")
	}
	// 密码不为空则更新（BeforeUpdate自动加密）
	if newPwd != "" {
		member.Password = newPwd
	}
	member.Major = major
	return db.DB.Save(&member).Error
}

// ListMembers 查询所有用户
func ListMembers() ([]model.Member, error) {
	var members []model.Member
	// 只查询需要的字段
	return members, db.DB.Select("id, account, major, created_at").Find(&members).Error
}

// GetMemberByID 根据ID查询用户详情
func GetMemberByID(id uint) (*model.Member, error) {
	var member model.Member
	if err := db.DB.Select("id, account, major, created_at").First(&member, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &member, nil
}
func CreateDefaultAdmin() error {
	// 定义默认管理员信息
	defaultAdminAccount := "admin"
	defaultAdminPwd := "Admin123!" // 建议设置复杂密码
	defaultAdminMajor := "软件工程"

	// 检查默认管理员是否已存在
	var exist model.Member
	if err := db.DB.Where("account = ?", defaultAdminAccount).First(&exist).Error; err == nil {
		// 已存在则直接返回成功
		return nil
	}
	// 不存在则创建
	return db.DB.Create(&model.Member{
		Account:   defaultAdminAccount,
		Password:  defaultAdminPwd,
		Major:     defaultAdminMajor,
		Character: "admin", // 直接设置为管理员角色
	}).Error

}
