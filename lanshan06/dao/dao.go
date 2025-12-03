package dao

import "golang.org/x/crypto/bcrypt"

// 模拟数据库
var database = map[string]string{}

func AddUser(username string, plainPassword string) error {
	//加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	database[username] = string(hashedPassword)
	return nil
}
func FindUser(username string, plainPassword string) bool {

	//从数据库获取密码
	hashedPassword, ok := database[username]
	if !ok {
		return false
	}
	//对比明文密码与机密密码是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
func SelectPasswordFromUsername(username string) string {
	return database[username]
}
func CheckUserExists(username string) bool {
	_, ok := database[username]
	return ok
}
func UpdatePassword(username string, newPlainPassword string) error {
	if !CheckUserExists(username) {
		return nil
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPlainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	database[username] = string(newHashedPassword)
	return nil
}
