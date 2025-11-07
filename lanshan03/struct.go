package main

import "fmt"

// 结构体嵌套********
type Info struct {
	name string
	age  int
}
type Poeple struct {
	Info
	group string
}

// ****************func(接收者变量名  接收者类型)方法名(参数列表)(返回参数)
func (p *Poeple) love() {
	fmt.Printf("%v:Dont forget what you are love", p.name)
}

// *******************接口**************************************
type sport interface {
	baseketball()
}
type ball struct{}

func (b ball) baseketball() {
	fmt.Println("hhh")
}
func main() {
	qsx := &Poeple{
		Info: Info{
			name: "邱双喜",
			age:  19,
		},
		group: "lllll",
	}

	fmt.Println(qsx.Info.name)
	fmt.Println(qsx.Info.age)
	fmt.Println(qsx.group)
	qsx.love()
	var sport sport
	sport = ball{}
	sport.baseketball()
}
