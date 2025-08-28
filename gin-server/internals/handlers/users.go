package handlers

import (
	"net/http"

	"github.com/krisn2/ginserver/internals/auth"
	"github.com/krisn2/ginserver/internals/db"

	// "github.com/krisn2/ginserver/internals/handlers"
	"github.com/krisn2/ginserver/internals/models"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	uid := auth.UID(c)
	var u models.User
	if err := db.DB.First(&u, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
}
