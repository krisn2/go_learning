package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/middleware"
	"github.com/krisn2/go-social/models"
	"gorm.io/gorm"
)

type LikeHandler struct {
	db *gorm.DB
}

func NewLikeHandler(db *gorm.DB) *LikeHandler {
	return &LikeHandler{db: db}
}

func (h *LikeHandler) Toggle(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var post models.Post
	if err := h.db.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Check if like exists
	var like models.Like
	err := h.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error

	if err == gorm.ErrRecordNotFound {
		// Create new like
		newLike := models.Like{UserID: userID, PostID: post.ID}
		if err := h.db.Create(&newLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"liked": true})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Delete existing like
	if err := h.db.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": false})
}
