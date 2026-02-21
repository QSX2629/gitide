package main

import (
	"fmt"
	"lanshan03/model"
)

// UserManager 用户管理器，负责用户的添加和管理
type UserManager struct {
	userList []model.User // 使用model包的User结构体
}

// Adder 定义添加用户的接口
type Adder interface {
	Add(user model.User) error
}

// Add 实现添加用户的方法，避免重复添加
func (u *UserManager) Add(user model.User) error {
	// 检查用户名是否已存在
	for _, existingUser := range u.userList {
		if existingUser.Name == user.Name {
			return fmt.Errorf("用户名:%v已存在", user.Name)
		}
	}
	u.userList = append(u.userList, user)
	return nil
}

func main() {
	manager := &UserManager{} // 修正拼写错误：Maneger → Manager

	// 准备用户数据（直接使用model.User）
	users := []model.User{
		{Name: "降魔大圣", Age: 2000, Gender: "male", Level: 90},
		{Name: "神里凝华", Age: 20, Gender: "female", Level: 90},
		{Name: "邱双喜", Age: 18, Gender: "male", Level: 59},
		{Name: "降魔大圣", Age: 200, Gender: "male", Level: 90}, // 重复用户
	}

	// 批量添加用户
	for _, user := range users {
		err := manager.Add(user)
		if err != nil {
			fmt.Printf("添加%s失败; %v\n", user.Name, err)
		} else {
			fmt.Printf("添加成功: %v\n", user.Name)
		}
	}
	fmt.Println()

	// 打印用户列表
	fmt.Println("用户列表:")
	for _, u := range manager.userList {
		fmt.Printf("%s（%d岁，%s，等级%d）\n", u.Name, u.Age, u.Gender, u.Level)
	}
}
