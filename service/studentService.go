package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type StudenController struct {
	db *gorm.DB
}

func NewStudenController(db *gorm.DB) *StudenController {
	return &StudenController{db: db}
}

// ==================================GetCourseInfo==================================

func (sc *StudenController) GetCourseInfo(c *gin.Context) {
	courseIdStr := c.Param("courseId")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var course models.Course
	if err := sc.db.Preload("Teacher").Preload("Students").First(&course, courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})

}

// ==================================GetCoursesForStudent==================================

func (sc *StudenController) GetCoursesForStudent(c *gin.Context) {
	studentIdStr := c.Param("studentId")
	studentId, err := strconv.ParseUint(studentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.User
	if err := sc.db.Preload("Courses").First(&student, studentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student.Courses})
}
