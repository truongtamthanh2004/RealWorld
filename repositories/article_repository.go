package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type ArticleRepository interface {
	//Create(article *models.Article) error
	//FindBySlug(slug string) (*models.Article, error)
	//Update(article *models.Article) error
	//Delete(article *models.Article) error
	List(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error)
	//Feed(userID uint, limit, offset int) ([]models.Article, int64, error)
}

type ArticleRepositoryImpl struct{}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{}
}

func (r *ArticleRepositoryImpl) List(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error) {
	var articles []models.Article
	query := database.DB.Preload("Author").Preload("Favorites").Model(&models.Article{})

	if tag != "" {
		query = query.Where("? = ANY (tag_list)", tag)
	}

	if author != "" {
		var user models.User
		if err := database.DB.Where("username = ?", author).First(&user).Error; err == nil {
			query = query.Where("author_id = ?", user.ID)
		} else {
			return []models.Article{}, 0, nil
		}
	}

	if favorited != "" {
		var favUser models.User
		if err := database.DB.Where("username = ?", favorited).First(&favUser).Error; err == nil {
			var favoriteArticleIDs []uint
			database.DB.Model(&models.Favorite{}).Where("user_id = ?", favUser.ID).Pluck("article_id", &favoriteArticleIDs)
			if len(favoriteArticleIDs) > 0 {
				query = query.Where("id IN ?", favoriteArticleIDs)
			} else {
				return []models.Article{}, 0, nil
			}
		} else {
			return []models.Article{}, 0, nil
		}
	}

	var totalCount int64
	query.Count(&totalCount)

	if limit == 0 {
		limit = 20
	}
	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	err := query.Find(&articles).Error
	return articles, totalCount, err
}
