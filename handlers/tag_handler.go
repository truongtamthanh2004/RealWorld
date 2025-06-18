package handlers

import (
	"Netlfy/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TagHandler struct {
	service services.TagService
}

func NewTagHandler(service services.TagService) *TagHandler {
	return &TagHandler{service}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.service.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
