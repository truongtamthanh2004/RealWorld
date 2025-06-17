package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Slug        string `gorm:"type:varchar(255);uniqueIndex"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	Body        string `gorm:"type:longtext"`
	AuthorID    uint
	Author      User       `gorm:"foreignKey:AuthorID"`
	TagList     []Tag      `gorm:"many2many:article_tags;"`
	Favorites   []Favorite `gorm:"foreignKey:ArticleID"`
	Comments    []Comment  `gorm:"foreignKey:ArticleID"`
}
