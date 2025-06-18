package services

import (
	"Netlfy/database"
	"Netlfy/dto"
	"Netlfy/models"
	"Netlfy/repositories"
	"Netlfy/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleService interface {
	ListArticles(tag, author, favorited string, limit, offset int, currentUserID *uint) ([]gin.H, int64, error)
	FeedArticles(userID uint, limit, offset int) ([]gin.H, int64, error)
	GetArticle(slug string, currentUserID *uint) (gin.H, error)
	CreateArticle(input dto.CreateArticleRequest, authorID uint) (map[string]interface{}, error)
	UpdateArticle(slug string, updatedData *dto.UpdateArticleRequest, userID uint) (gin.H, error)
	DeleteArticle(slug string, userID uint) error
	FavoriteArticle(slug string, userID uint) (gin.H, error)
	UnfavoriteArticle(slug string, userID uint) (gin.H, error)
}

type ArticleServiceImpl struct {
	articleRepo  repositories.ArticleRepository
	followRepo   repositories.FollowRepository
	favoriteRepo repositories.FavoriteRepository
}

func NewArticleService(articleRepo repositories.ArticleRepository, followRepo repositories.FollowRepository, favoriteRepo repositories.FavoriteRepository) ArticleService {
	return &ArticleServiceImpl{articleRepo, followRepo, favoriteRepo}
}

func (s *ArticleServiceImpl) ListArticles(tag, author, favorited string, limit, offset int, currentUserID *uint) ([]gin.H, int64, error) {
	articles, count, err := s.articleRepo.List(tag, author, favorited, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var serializedArticles []gin.H
	for _, article := range articles {
		serializedArticles = append(serializedArticles, s.SerializeArticle(&article, currentUserID))
	}
	return serializedArticles, count, nil
}

func (s *ArticleServiceImpl) FeedArticles(userID uint, limit, offset int) ([]gin.H, int64, error) {
	articles, count, err := s.articleRepo.Feed(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var serializedArticles []gin.H
	for _, article := range articles {
		serializedArticles = append(serializedArticles, s.SerializeArticle(&article, &userID))
	}
	return serializedArticles, count, nil
}

func (s *ArticleServiceImpl) GetArticle(slug string, currentUserID *uint) (gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	return s.SerializeArticle(article, currentUserID), nil
}

func (s *ArticleServiceImpl) CreateArticle(input dto.CreateArticleRequest, authorID uint) (map[string]interface{}, error) {
	var tags []models.Tag
	for _, tagName := range input.TagList {
		var tag models.Tag
		if err := database.DB.Where("name = ?", tagName).First(&tag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tag = models.Tag{Name: tagName}
				if err := database.DB.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		tags = append(tags, tag)
	}

	newArticle := &models.Article{
		Slug:        utils.GenerateSlug(input.Title),
		Title:       input.Title,
		Description: input.Description,
		Body:        input.Body,
		TagList:     tags,
		AuthorID:    authorID,
	}

	if err := s.articleRepo.Create(newArticle); err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Author").First(&newArticle, newArticle.ID).Error; err != nil {
		return nil, err
	}

	return s.SerializeArticle(newArticle, &authorID), nil
}

func (s *ArticleServiceImpl) UpdateArticle(slug string, updatedData *dto.UpdateArticleRequest, userID uint) (gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	if article.AuthorID != userID {
		return nil, errors.New("unauthorized to update this article")
	}

	if updatedData.Title != "" {
		article.Title = updatedData.Title
		article.Slug = utils.GenerateSlug(updatedData.Title)
	}
	if updatedData.Description != "" {
		article.Description = updatedData.Description
	}
	if updatedData.Body != "" {
		article.Body = updatedData.Body
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return s.SerializeArticle(article, &userID), nil
}

func (s *ArticleServiceImpl) DeleteArticle(slug string, userID uint) error {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return err
	}

	if article.AuthorID != userID {
		return errors.New("unauthorized to delete this article")
	}

	return s.articleRepo.DeleteByID(article.ID)
}

func (s *ArticleServiceImpl) FavoriteArticle(slug string, userID uint) (gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	err = s.favoriteRepo.AddFavorite(userID, article.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{"article": s.SerializeArticle(article, &userID)}, nil
}

func (s *ArticleServiceImpl) UnfavoriteArticle(slug string, userID uint) (gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	err = s.favoriteRepo.RemoveFavorite(userID, article.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{"article": s.SerializeArticle(article, &userID)}, nil
}

func (s *ArticleServiceImpl) SerializeArticle(article *models.Article, currentUserID *uint) gin.H {
	var favorited bool
	var following bool
	if currentUserID != nil {
		var fav models.Favorite
		err := database.DB.Where("user_id = ? AND article_id = ?", *currentUserID, article.ID).First(&fav).Error
		favorited = err == nil

		following, _ = s.followRepo.IsFollowing(*currentUserID, article.AuthorID)
	}

	return gin.H{
		"slug":           article.Slug,
		"title":          article.Title,
		"description":    article.Description,
		"body":           article.Body,
		"tagList":        article.TagList,
		"createdAt":      article.CreatedAt,
		"updatedAt":      article.UpdatedAt,
		"favorited":      favorited,
		"favoritesCount": len(article.Favorites),
		"author": gin.H{
			"username":  article.Author.Username,
			"bio":       article.Author.Bio,
			"image":     article.Author.Image,
			"following": following,
		},
	}
}
