package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name     string     `gorm:"unique"`
	Articles []*Article `gorm:"many2many:article_tags;"`
}
