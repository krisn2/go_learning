package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/models"
)

func Paginate(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	if v := c.Query("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}

func PostResponse(post models.Post, likesCount, commentsCount int) gin.H {
	return gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"body":       post.Body,
		"author":     UserResponse(post.User),
		"likes":      likesCount,
		"comments":   commentsCount,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	}
}

func UserResponse(user models.User) gin.H {
	return gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
}

func CommentResponse(comment models.Comment) gin.H {
	return gin.H{
		"id":         comment.ID,
		"body":       comment.Body,
		"author":     UserResponse(comment.User),
		"created_at": comment.CreatedAt,
	}
}
