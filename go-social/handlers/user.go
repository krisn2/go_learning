package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/config"
	"github.com/krisn2/go-social/middleware"
	"github.com/krisn2/go-social/models"
	"github.com/krisn2/go-social/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	response := utils.UserResponse(user)
	response["created_at"] = user.CreatedAt
	response["updated_at"] = user.UpdatedAt

	c.JSON(http.StatusOK, response)
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"omitempty,min=2,max=100"`
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Update("name", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	h.GetMe(c)
}

func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Delete user's comments and likes
		if err := tx.Where("user_id = ?", userID).Delete(&models.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Like{}).Error; err != nil {
			return err
		}

		// Get user's posts
		var posts []models.Post
		if err := tx.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
			return err
		}

		// Delete comments and likes on user's posts
		if len(posts) > 0 {
			var postIDs []uint
			for _, post := range posts {
				postIDs = append(postIDs, post.ID)
			}

			if err := tx.Where("post_id IN ?", postIDs).Delete(&models.Comment{}).Error; err != nil {
				return err
			}
			if err := tx.Where("post_id IN ?", postIDs).Delete(&models.Like{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", postIDs).Delete(&models.Post{}).Error; err != nil {
				return err
			}
		}

		// Delete user
		return tx.Delete(&models.User{}, userID).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// Demo admin check - replace with proper RBAC in production
	if c.Query("admin_secret") != h.cfg.AdminListSecret {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	page, pageSize := utils.Paginate(c)

	var users []models.User
	var total int64

	h.db.Model(&models.User{}).Count(&total)
	h.db.Select("id, name, email, created_at, updated_at").
		Order("id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
