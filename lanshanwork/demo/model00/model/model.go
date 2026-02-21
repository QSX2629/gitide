package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Member 会员模型（数据库表映射）
type Member struct {
	gorm.Model
	Account   string `gorm:"not null;unique;comment:账号"`
	Password  string `gorm:"not null;comment:密码（加密后存储）"`
	Major     string `gorm:"not null;comment:专业"`
	Character string `gorm:"not null;comment:"角色""`
}

// BeforeCreate 创建前加密密码
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hash)
	return nil
}

// BeforeUpdate 更新前加密新密码（仅当密码修改时）
func (m *Member) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") { // 仅密码字段修改时加密
		hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		m.Password = string(hash)
	}
	return nil
}

// CheckPassword 校验密码
func (m *Member) CheckPassword(plainPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(plainPwd)) == nil
}
