package services

import (
	"Netlfy/database"
	"Netlfy/models"
	"Netlfy/repositories"
	"errors"
	"github.com/gin-gonic/gin"
)

type CommentService interface {
	AddComment(slug string, body string, userID uint) (gin.H, error)
	GetComments(slug string, currentUserID *uint) ([]gin.H, error)
	DeleteComment(slug string, commentID uint, userID uint) error
}

type CommentServiceImpl struct {
	commentRepo repositories.CommentRepository
	articleRepo repositories.ArticleRepository
	followRepo  repositories.FollowRepository
}

func NewCommentService(commentRepo repositories.CommentRepository, articleRepo repositories.ArticleRepository, followRepo repositories.FollowRepository) CommentService {
	return &CommentServiceImpl{commentRepo, articleRepo, followRepo}
}

func (s *CommentServiceImpl) AddComment(slug string, body string, userID uint) (gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	comment := models.Comment{
		Body:      body,
		ArticleID: article.ID,
		AuthorID:  userID,
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	// Load Author for serialization
	database.DB.Preload("Author").First(&comment, comment.ID)

	return gin.H{"comment": s.SerializeComment(&comment, &userID)}, nil
}

func (s *CommentServiceImpl) GetComments(slug string, currentUserID *uint) ([]gin.H, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	var comments []models.Comment
	if err := database.DB.Preload("Author").Where("article_id = ?", article.ID).Find(&comments).Error; err != nil {
		return nil, err
	}

	var result []gin.H
	for _, comment := range comments {
		result = append(result, s.SerializeComment(&comment, currentUserID))
	}
	return result, nil
}

func (s *CommentServiceImpl) DeleteComment(slug string, commentID uint, userID uint) error {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return err
	}

	var comment models.Comment
	err = database.DB.Where("id = ? AND article_id = ?", commentID, article.ID).First(&comment).Error
	if err != nil {
		return err
	}

	// Only author can delete
	if comment.AuthorID != userID {
		return errors.New("forbidden")
	}

	return database.DB.Delete(&comment).Error
}

func (s *CommentServiceImpl) SerializeComment(comment *models.Comment, currentUserID *uint) gin.H {
	var following bool
	if currentUserID != nil {
		following, _ = s.followRepo.IsFollowing(*currentUserID, comment.AuthorID)
	}

	return gin.H{
		"id":        comment.ID,
		"createdAt": comment.CreatedAt,
		"updatedAt": comment.UpdatedAt,
		"body":      comment.Body,
		"author": gin.H{
			"username":  comment.Author.Username,
			"bio":       comment.Author.Bio,
			"image":     comment.Author.Image,
			"following": following,
		},
	}
}
