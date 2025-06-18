package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByArticleSlug(slug string) ([]models.Comment, error)
	Delete(commentID uint, authorID uint) error
}

type CommentRepositoryImpl struct{}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}

func (r *CommentRepositoryImpl) Create(comment *models.Comment) error {
	return database.DB.Create(comment).Error
}

func (r *CommentRepositoryImpl) FindByArticleSlug(slug string) ([]models.Comment, error) {
	var article models.Article
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}

	var comments []models.Comment
	err := database.DB.Where("article_id = ?", article.ID).Preload("Author").Order("created_at").Find(&comments).Error
	return comments, err
}

func (r *CommentRepositoryImpl) Delete(commentID uint, authorID uint) error {
	return database.DB.Where("id = ? AND author_id = ?", commentID, authorID).Delete(&models.Comment{}).Error
}
