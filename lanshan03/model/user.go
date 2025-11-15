package model

// User 定义用户数据结构
type User struct {
	Name   string // 首字母大写，允许跨包访问
	Age    int
	Gender string
	Level  int
}
