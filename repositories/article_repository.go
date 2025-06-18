package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type ArticleRepository interface {
	Create(article *models.Article) error
	FindBySlug(slug string) (*models.Article, error)
	Update(article *models.Article) error
	DeleteByID(articleID uint) error
	List(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error)
	Feed(userID uint, limit, offset int) ([]models.Article, int64, error)
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

func (r *ArticleRepositoryImpl) Feed(userID uint, limit, offset int) ([]models.Article, int64, error) {
	var followees []uint
	if err := database.DB.Model(&models.Follow{}).Where("follower_id = ?", userID).Pluck("followee_id", &followees).Error; err != nil {
		return nil, 0, err
	}
	if len(followees) == 0 {
		return []models.Article{}, 0, nil
	}

	var articles []models.Article
	query := database.DB.Preload("Author").Preload("Favorites").Where("author_id IN ?", followees)
	var totalCount int64
	query.Model(&models.Article{}).Count(&totalCount)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, totalCount, err
}

func (r *ArticleRepositoryImpl) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := database.DB.Preload("Author").Preload("Favorites").Where("slug = ?", slug).First(&article).Error
	return &article, err
}

func (r *ArticleRepositoryImpl) Create(article *models.Article) error {
	return database.DB.Create(article).Error
}

func (r *ArticleRepositoryImpl) Update(article *models.Article) error {
	return database.DB.Save(article).Error
}

func (r *ArticleRepositoryImpl) DeleteByID(articleID uint) error {
	return database.DB.Delete(&models.Article{}, articleID).Error
}
