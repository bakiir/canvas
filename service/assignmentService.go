package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type AssignmentController struct {
	db *gorm.DB
}

func NewAssignmentController(db *gorm.DB) *AssignmentController {
	return &AssignmentController{db: db}
}

func (tc *TeacherController) GetAssignmentsByCourse(c *gin.Context) {
	courseIdStr := c.Param("courseId")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var assignments []models.Assignment
	if err := tc.db.Where("course_id = ?", courseId).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignments})
}
