package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetAllCourses(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var courses []models.Course
	if err := db.Preload("Teacher").Preload("Students").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}
