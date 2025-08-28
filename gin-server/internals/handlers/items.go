package handlers

import (
	"net/http"
	"strconv"

	"github.com/krisn2/ginserver/internals/auth"
	"github.com/krisn2/ginserver/internals/db"
	"github.com/krisn2/ginserver/internals/models"

	"github.com/gin-gonic/gin"
)

type itemDTO struct {
	Title string `json:"title" binding:"required"`
}

// POST /api/items
func CreateItem(c *gin.Context) {
	uid := auth.UID(c)
	var in itemDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it := models.Item{Title: in.Title, UserID: uid}
	db.DB.Create(&it)
	c.JSON(http.StatusCreated, it)
}

// GET /api/items/:id  (path param)
func GetItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var it models.Item
	if err := db.DB.First(&it, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, it)
}

// GET /api/items?limit=10  (query param)
func ListItems(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var items []models.Item
	db.DB.Limit(limit).Order("id desc").Find(&items)
	c.JSON(http.StatusOK, items)
}
