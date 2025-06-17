package services

import (
	"Netlfy/database"
	"Netlfy/models"
	"Netlfy/repositories"
	"github.com/gin-gonic/gin"
)

type ArticleService interface {
	ListArticles(tag, author, favorited string, limit, offset int, currentUserID *uint) ([]gin.H, int64, error)
}

type ArticleServiceImpl struct {
	articleRepo repositories.ArticleRepository
	followRepo  repositories.FollowRepository
}

func NewArticleService(articleRepo repositories.ArticleRepository, followRepo repositories.FollowRepository) ArticleService {
	return &ArticleServiceImpl{articleRepo, followRepo}
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
