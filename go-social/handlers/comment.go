package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/middleware"
	"github.com/krisn2/go-social/models"
	"github.com/krisn2/go-social/utils"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

type CommentRequest struct {
	Body string `json:"body" binding:"required,min=1"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var post models.Post
	if err := h.db.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		Body:   req.Body,
		UserID: userID,
		PostID: post.ID,
	}

	if err := h.db.Create(comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	if err := h.db.Preload("User").First(comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusCreated, gin.H{"id": comment.ID})
		return
	}

	c.JSON(http.StatusCreated, utils.CommentResponse(*comment))
}

func (h *CommentHandler) List(c *gin.Context) {
	page, pageSize := utils.Paginate(c)

	var comments []models.Comment
	if err := h.db.Preload("User").
		Where("post_id = ?", c.Param("id")).
		Order("id ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}

	response := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		response = append(response, utils.CommentResponse(comment))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var comment models.Comment
	if err := h.db.First(&comment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	var post models.Post
	h.db.Select("id, user_id").First(&post, comment.PostID)

	// Check if user is comment author or post author
	if comment.UserID != userID && post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to delete this comment"})
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
