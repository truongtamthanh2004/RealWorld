package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Body      string
	ArticleID uint
	AuthorID  uint
	Author    User
}
