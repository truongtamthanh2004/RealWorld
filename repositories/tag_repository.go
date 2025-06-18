package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type TagRepository interface {
	GetAllTags() ([]models.Tag, error)
}

type TagRepositoryImpl struct{}

func NewTagRepository() TagRepository {
	return &TagRepositoryImpl{}
}

func (r *TagRepositoryImpl) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := database.DB.Find(&tags).Error
	return tags, err
}
