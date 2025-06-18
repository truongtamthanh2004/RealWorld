package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type FavoriteRepository interface {
	AddFavorite(userID, articleID uint) error
	RemoveFavorite(userID, articleID uint) error
	IsFavorited(userID, articleID uint) (bool, error)
	CountFavorites(articleID uint) (int64, error)
}

type FavoriteRepositoryImpl struct{}

func NewFavoriteRepository() FavoriteRepository {
	return &FavoriteRepositoryImpl{}
}

func (r *FavoriteRepositoryImpl) AddFavorite(userID, articleID uint) error {
	return database.DB.Create(&models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}).Error
}

func (r *FavoriteRepositoryImpl) RemoveFavorite(userID, articleID uint) error {
	return database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{}).Error
}

func (r *FavoriteRepositoryImpl) IsFavorited(userID, articleID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepositoryImpl) CountFavorites(articleID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Favorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
