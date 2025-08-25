package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/middleware"
	"github.com/krisn2/go-social/models"
	"github.com/krisn2/go-social/utils"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

type PostRequest struct {
	Title string `json:"title" binding:"required,min=1,max=200"`
	Body  string `json:"body" binding:"max=10000"` // Add max length
}

func (h *PostHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &models.Post{
		Title:  req.Title,
		Body:   req.Body,
		UserID: userID,
	}

	if err := h.db.Create(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	// Load user for response
	h.db.Preload("User").First(post, post.ID)
	h.getPostResponse(c, post)
}

func (h *PostHandler) List(c *gin.Context) {
	page, pageSize := utils.Paginate(c)

	// Optimized query with subqueries for counts
	var posts []struct {
		models.Post
		LikesCount    int64 `json:"likes_count"`
		CommentsCount int64 `json:"comments_count"`
	}

	var total int64
	h.db.Model(&models.Post{}).Count(&total)

	// Single query with counts
	h.db.Table("posts").
		Select(`posts.*, 
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as likes_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comments_count`).
		Joins("LEFT JOIN users ON posts.user_id = users.id").
		Order("posts.id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&posts)

	// Load users separately to avoid N+1
	var userIDs []uint
	for _, post := range posts {
		userIDs = append(userIDs, post.UserID)
	}

	var users []models.User
	h.db.Where("id IN ?", userIDs).Find(&users)

	// Create user map for O(1) lookup
	userMap := make(map[uint]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	response := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		post.User = userMap[post.UserID]
		response = append(response, utils.PostResponse(post.Post, int(post.LikesCount), int(post.CommentsCount)))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *PostHandler) Get(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var post models.Post
	if err := h.db.Preload("User").First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	h.getPostResponse(c, &post)
}

func (h *PostHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var post models.Post
	if err := h.db.First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the post owner"})
		return
	}

	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&post).Updates(models.Post{Title: req.Title, Body: req.Body}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	h.db.Preload("User").First(&post, post.ID)
	h.getPostResponse(c, &post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var post models.Post
	if err := h.db.First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not the post owner"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("post_id = ?", post.ID).Delete(&models.Like{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Post{}, post.ID).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *PostHandler) getPostResponse(c *gin.Context, post *models.Post) {
	likesCount, commentsCount := h.getCounts(post.ID)
	c.JSON(http.StatusOK, utils.PostResponse(*post, likesCount, commentsCount))
}

func (h *PostHandler) getCounts(postID uint) (int, int) {
	var likesCount int64
	var commentsCount int64

	h.db.Model(&models.Like{}).Where("post_id = ?", postID).Count(&likesCount)
	h.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&commentsCount)

	return int(likesCount), int(commentsCount)
}
