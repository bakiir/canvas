package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type TeacherController struct {
	db *gorm.DB
}

func NewTeacherController(db *gorm.DB) *TeacherController {
	return &TeacherController{db: db}
}

// ==================================GetCoursesForTeacher==================================

func (sc *TeacherController) GetCoursesForTeacher(c *gin.Context) {
	teacherIdStr := c.Param("teacherId")
	teacherId, err := strconv.ParseUint(teacherIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var teacher models.User
	if err := sc.db.Preload("Courses").First(&teacher, teacherId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teacher.Courses})
}

// ==================================AddAssignmentToCourse==================================

func (tc *TeacherController) AddAssignmentToCourse(c *gin.Context) {
	courseIdStr := c.Param("courseId")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, что курс существует
	var course models.Course
	if err := tc.db.First(&course, courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Назначаем ID курса заданию
	input.CourseID = uint(courseId)

	if err := tc.db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// ==================================RemoveAssignmentFromCourse==================================

func (tc *TeacherController) RemoveAssignmentFromCourse(c *gin.Context) {
	assignmentIdStr := c.Param("assignmentId")
	assignmentId, err := strconv.ParseUint(assignmentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var input models.Assignment
	if err := tc.db.First(&input, assignmentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if err := tc.db.Delete(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove assignment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}
