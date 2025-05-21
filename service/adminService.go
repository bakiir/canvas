package service

import (
	"CanvasApplication/models"
	"CanvasApplication/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type AdminController struct {
	db *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{db: db}
}

func (ac *AdminController) CreateStudent(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности логина
	var existing models.Student
	if err := ac.db.Where("login = ?", input.Login).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
		return
	}

	student := models.Student{
		Name:     input.Name,
		Login:    input.Login,
		Password: input.Password, // В реальном приложении хешируйте пароль!
	}

	if err := ac.db.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	// Не возвращаем пароль в ответе
	student.Password = ""
	c.JSON(http.StatusCreated, gin.H{"data": student})
}

func (ac *AdminController) CreateTeacher(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности логина
	var existing models.Teacher
	if err := ac.db.Where("login = ?", input.Login).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
		return
	}

	teacher := models.Teacher{
		Name:     input.Name,
		Login:    input.Login,
		Password: input.Password, // В реальном приложении хешируйте пароль!
	}

	if err := ac.db.Create(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher"})
		return
	}

	// Не возвращаем пароль в ответе
	teacher.Password = ""
	c.JSON(http.StatusCreated, gin.H{"data": teacher})
}

func (ac *AdminController) DeleteStudent(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := ac.db.Where("id = ?", studentID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	ac.db.Delete(&student) // Полное удаление
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func (ac *AdminController) DeleteTeacher(c *gin.Context) {
	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}
	var teacher models.Teacher
	if err := ac.db.Where("id = ?", teacherID).First(&teacher).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
	}

	c.JSON(http.StatusOK, gin.H{"deleted": teacher})

}

func (ac *AdminController) DeleteCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	var course models.Course
	if err := ac.db.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})

	}
	ac.db.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})

}

func (tc *AdminController) Login(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := tc.db.Where("login = ?", credentials.Login).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !student.CheckPassword(credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(student.ID, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (tc *AdminController) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Student
	if err := tc.db.Where("login = ?", input.Login).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
		return
	}

	student := models.Student{
		Name:  input.Name,
		Login: input.Login,
	}

	if err := student.SetPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := tc.db.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully"})
}
