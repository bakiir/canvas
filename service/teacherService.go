package service

import (
	"CanvasApplication/models"
	"CanvasApplication/utils"
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

	var existing models.Teacher
	if err := tc.db.Where("login = ?", input.Login).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
		return
	}

	teacher := models.Teacher{
		Name:  input.Name,
		Login: input.Login,
	}

	if err := teacher.SetPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := tc.db.Create(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Teacher created successfully"})
}

// Логин преподавателя
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

	if !teacher.CheckPassword(credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(teacher.ID, "teacher")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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

	var teacher models.Teacher
	if err := tc.db.Preload("Courses").First(&teacher, input.TeacherID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

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
