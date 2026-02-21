package model

import (
	"gorm.io/gorm"
)

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusDeleted   = "deleted"
)

type User struct {
	gorm.Model
	Account     string `gorm:"size:20;DEFAULT:'0'"json:"account"`
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `gorm:"size:255;uniqueIndex;fulltext" json:"name"`
	Password    string `gorm:"size:255" json:"password"`
	Common_user bool   `gorm:"not null;default:false" json:"common_user"`
	Admin_user  bool   `gorm:"not null;default:false" json:"admin_user"`
}
type Article struct {
	gorm.Model
	Title     string `gorm:"size:255;uniqueIndex;fulltext" json:"title"`
	Content   string `gorm:"size:255" json:"content"`
	ArticleID uint   `gorm: json:"id"`
	UserId    uint   `gorm:"index" json:"user_id"`
	Status    string `gorm:"size:255" json:"status;default:'draft'"`
}
type Comment struct {
	gorm.Model
	CommentID uint   `gorm: json:"id"`
	Content   string `gorm:"size:255" json:"content"`
	ArticleId uint   `gorm:"index" json:"article_id"`
}
type Follow struct {
	gorm.Model
	ID         uint `gorm:"primary" json:"id"`
	UserID     uint `gorm:"not null;index:idx_user_followed" json:"user_id"` // 关注者ID
	FollowedID uint `gorm:"not null;index:idx_user_followed" json:"followed_id"`
} // 被关注者ID
