package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type CourseController struct {
	db *gorm.DB
}

func NewCourseController(db *gorm.DB) *CourseController {
	return &CourseController{db: db}
}

// ==================================GetCoursesByTeacher==================================

func (cc *CourseController) GetCoursesByTeacher(c *gin.Context) {
	teacherIdStr := c.Param("teacher_id")
	teacherId, err := strconv.ParseUint(teacherIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var courses []models.Course
	if err := cc.db.Where("teacher_id", teacherId).Find(&courses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": courses})
}

// ==================================GetCoursesByStudent==================================

func (cc *CourseController) GetCoursesByStudent(c *gin.Context) {
	studentIdStr := c.Param("studentId")
	studentId, err := strconv.ParseUint(studentIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var courses []models.Course

	if err := cc.db.Where("student_id", studentId).Find(&courses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": courses})
}

// ==================================DeleteCourse==================================

func (cc *CourseController) DeleteCourse(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := cc.db.Where("course_id", courseId).Delete(&models.Course{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Deleted"})
}
