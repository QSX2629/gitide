package model01

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"gorm:"unique;not null"`
	Body    string `json:"body"`
	Account string `json:"account"` //对应相应的发布账户
}
type Comment struct {
	gorm.Model
	Content     string `json:"content"`
	ArticleName string `json:"article_name"` //对应唯一的文章名字
}
