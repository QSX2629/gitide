package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username string `json:"username"binding:required`
	Password string `json:"password"binding:required,min=6,ma=20`
}

// HashPassword 密码加密：明文的密码转为哈希值储存到数据库
func (u *User) HashPassword() error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedBytes)
	return nil
}

// CheckPassword 密码校验：明文密码与哈希值是否一致
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

type ModifyPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20"`       // 旧密码：必填，长度6-20
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`       // 新密码：必填，长度6-20
	ConfirmPwd  string `json:"confirm_pwd" binding:"required,eqfield=NewPassword"` // 确认新密码：必填，需与新密码一致
}
