package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Member 定义会员数据模型
type Member struct {
	Account  string `gorm:"not null;unique;comment:"账号""`
	Password string `gorm:"not null;min:6;max:10;comment:"密码""`
	Major    string `gorm:"not null;comment:专业"`
}

// BeforeCreate ////////////////*****************构建数据模型********************/////////////////////
// BeforeCreate GORM 钩子：在创建 Member 前加密密码
func (m *Member) BeforeCreate(tx *gorm.DB) (err error) {
	// 对密码进行 bcrypt 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err // 加密失败则返回错误
	}
	m.Password = string(hashedPassword) // 替换原始密码为加密后的哈希值
	return nil
}
