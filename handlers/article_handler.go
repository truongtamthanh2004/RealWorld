package handlers

import (
	"Netlfy/dto"
	"Netlfy/services"
	"Netlfy/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService}
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var currentUserID *uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if userID, err := utils.ParseToken(tokenString); err == nil {
			currentUserID = &userID
		}
	}

	articles, count, err := h.articleService.ListArticles(tag, author, favorited, limit, offset, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":      articles,
		"articlesCount": count,
	})
}

func (h *ArticleHandler) FeedArticles(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	currentUserID := userIDValue.(uint)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	articles, count, err := h.articleService.FeedArticles(currentUserID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":      articles,
		"articlesCount": count,
	})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	var currentUserID *uint
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if userID, err := utils.ParseToken(token); err == nil {
			currentUserID = &userID
		}
	}

	article, err := h.articleService.GetArticle(slug, currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req struct {
		Article dto.CreateArticleRequest `json:"article"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(uint)

	createdArticle, err := h.articleService.CreateArticle(req.Article, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": createdArticle})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	var req struct {
		Article dto.UpdateArticleRequest `json:"article"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	slug := c.Param("slug")

	userIDValue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	result, err := h.articleService.UpdateArticle(slug, &req.Article, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": result})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")

	userIDValue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	err := h.articleService.DeleteArticle(slug, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ArticleHandler) FavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	user, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(uint)

	result, err := h.articleService.FavoriteArticle(slug, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ArticleHandler) UnfavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	user, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := user.(uint)

	result, err := h.articleService.UnfavoriteArticle(slug, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
