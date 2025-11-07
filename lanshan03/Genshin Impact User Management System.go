package main

import (
	"fmt"
)

type user struct {
	name   string
	age    int
	gender string
	level  int
}
type UserManeger struct {
	userlist []user
}

type userConf struct {
	name   string
	age    int
	gender string
	level  int
}

type Adder interface {
	Add(conf userConf) error
}

func (u *UserManeger) Add(conf userConf) error {
	Newuser := user{
		name:   conf.name,
		age:    conf.age,
		gender: conf.gender,
		level:  conf.level,
	}

	for _, user := range u.userlist {
		if user.name == Newuser.name {
			return fmt.Errorf("  用户名:%v已存在", Newuser.name)
		}
	}
	u.userlist = append(u.userlist, Newuser)
	return nil
}
func main() {
	maneger := &UserManeger{}
	confs := []userConf{
		{name: "降魔大圣", age: 2000, gender: "male", level: 90},
		{name: "神里凝华", age: 20, gender: "female", level: 90},
		{name: "邱双喜", age: 18, gender: "male", level: 59},
		{name: "降魔大圣", age: 200, gender: "male", level: 90},
	}

	for _, conf := range confs {
		err := maneger.Add(conf)
		if err != nil {
			fmt.Printf("添加%s失败;%v\n", conf.name, err)
		} else {
			fmt.Printf("添加成功%v\n", conf.name)
		}
	}
	fmt.Println()

	fmt.Println("用户列表:\n")
	for _, u := range maneger.userlist {
		fmt.Printf("%s（%d岁，%s，等级%d）\n", u.name, u.age, u.gender, u.level)
	}
}
