package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type StudentController struct {
	db *gorm.DB
}

func NewStudentController(db *gorm.DB) *StudentController {
	return &StudentController{db: db}
}

func (sc *StudentController) GetStudentCourses(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := sc.db.Preload("Courses").First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student.Courses})
}

func (sc *StudentController) GetAllStudents(c *gin.Context) {
	var students []models.Student
	if err := sc.db.Preload("Courses").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": students})
}

func (sc *StudentController) Enroll(c *gin.Context) {
	var input struct {
		StudentID uint `json:"student_id" binding:"required"`
		CourseID  uint `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование студента
	var student models.Student
	if err := sc.db.First(&student, input.StudentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Проверяем существование курса
	var course models.Course
	if err := sc.db.Preload("Students").First(&course, input.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Проверяем есть ли место на курсе
	if len(course.Students) >= course.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course is full"})
		return
	}

	// Проверяем уже записан ли студент
	count := sc.db.Model(&student).Association("Courses").Count()
	if count >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student has reached course limit (3)"})
		return
	}

	// Записываем студента на курс
	if err := sc.db.Model(&student).Association("Courses").Append(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll"})
		return
	}

	// Уменьшаем доступные места
	if err := sc.db.Model(&course).Update("capacity", course.Capacity-1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update capacity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Enrolled successfully",
		"capacity": course.Capacity - 1, // Возвращаем обновленную вместимость
	})
}

func (sc *StudentController) Login(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := sc.db.Where("login = ?", credentials.Login).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// В реальном приложении используйте bcrypt или аналоги
	if student.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Не возвращаем пароль
	student.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": student})
}
