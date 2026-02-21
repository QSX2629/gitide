package service

import (
	"context"
	"demo/db"
	"demo/model00/model"
	"demo/rediss"
	"demo/utils_lock"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Register 用户注册
func Register(account, password, major string) error {
	ctx := context.Background()
	// 1. 定义分布式锁的key：register:账号（唯一标识）
	lockKey := fmt.Sprintf("register:%s", account)
	// 2. 尝试获取锁（过期时间5秒，防止死锁）
	lockValue, ok, err := utils_lock.Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return errors.New("获取锁失败：" + err.Error())
	}
	if !ok {
		return errors.New("当前账号正在注册中，请稍后再试")
	}
	// 3. 延迟释放锁（确保函数执行完后释放）
	defer func() {
		_ = utils_lock.Unlock(ctx, lockKey, lockValue)
	}()
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
	//删除缓存
	cacheKey := fmt.Sprintf("member:%d", id)
	_ = rediss.DeleteCache(context.Background(), cacheKey)
	return nil
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
	if err := db.DB.Save(&member).Error; err != nil {
		return errors.New("更新失败")
	}
	//删除缓存（防止缓存与数据库数据不同）
	cacheKey := fmt.Sprintf("member:%d", id)
	_ = rediss.DeleteCache(context.Background(), cacheKey)
	return nil
}

// ListMembers 查询所有用户
func ListMembers() ([]model.Member, error) {
	var members []model.Member
	// 只查询需要的字段
	return members, db.DB.Select("id, account, major, created_at").Find(&members).Error
}

// GetMemberByID 根据ID查询用户详情
func GetMemberByID(id uint) (*model.Member, error) {
	//1.定义缓存的key
	cacheKey := fmt.Sprintf("%d", id)
	ctx := context.Background()
	//2.获取缓存
	cacheData, err := rediss.GetCache(ctx, cacheKey)
	if err == nil {
		var member model.Member
		//缓存的JSON解析为member
		if err = json.Unmarshal([]byte(cacheData), &member); err != nil {
			return nil, errors.New("缓存解析失败")
		}
		return &member, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, errors.New("查询失败" + err.Error())
	}
	//若没有缓存则访问数据库
	var member model.Member
	if err := db.DB.Select("id,account,major,created_at").First(&member, id).Error; err == nil {
		return nil, errors.New("用户不存在")
	}
	//将结果存入redis（设置过期时间为5分钟）
	memberJson, err := json.Marshal(member)
	if err != nil {
		return nil, errors.New("数据列表化失败")
	}
	_ = rediss.SetCache(ctx, cacheKey, string(memberJson), 60)
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
