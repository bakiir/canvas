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

func (tc *TeacherController) GetAllTeachers(c *gin.Context) {
	var teachers []models.Teacher
	if err := tc.db.Preload("Courses").Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}

	// Очищаем пароли перед отправкой
	for i := range teachers {
		teachers[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": teachers})
}

func (tc *TeacherController) TeacherLogin(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var teacher models.Teacher
	if err := tc.db.Where("login = ?", credentials.Login).First(&teacher).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if teacher.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	teacher.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": teacher})
}

func (tc *TeacherController) CreateCourse(c *gin.Context) {
	var input struct {
		TeacherID uint   `json:"teacher_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Capacity  int    `json:"capacity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование преподавателя
	var teacher models.Teacher
	if err := tc.db.Preload("Courses").First(&teacher, input.TeacherID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	// Получаем количество курсов преподавателя (правильный способ)
	count := tc.db.Model(&teacher).Association("Courses").Count()

	if count >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher has reached course limit (3)"})
		return
	}

	course := models.Course{
		Name:      input.Name,
		Capacity:  input.Capacity,
		TeacherID: input.TeacherID,
	}

	if err := tc.db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": course})
}

func (tc *TeacherController) GetTeacherCourses(c *gin.Context) {
	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID"})
		return
	}

	var teacher models.Teacher
	if err := tc.db.Preload("Courses").First(&teacher, teacherID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teacher.Courses})
}

func (tc *TeacherController) CreateTeacher(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем уникальность логина
	var existing models.Teacher
	if err := tc.db.Where("login = ?", input.Login).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
		return
	}

	teacher := models.Teacher{
		Name:     input.Name,
		Login:    input.Login,
		Password: input.Password, // В реальном приложении хешируйте!
	}

	if err := tc.db.Create(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher"})
		return
	}

	teacher.Password = ""
	c.JSON(http.StatusCreated, gin.H{"data": teacher})
}
