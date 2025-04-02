package service

import (
	"CanvasApplication/models"
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

//func (ac *AdminController) DeleteCourse(ctx *gin.Context) {
//	//studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//	//if err != nil {
//	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
//	//	return
//	//}
//	//var student models.Student
//	//if err := ac.db.Where("id = ?", studentID).First(&student).Error; err != nil {
//	//	c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
//	//}
//	//
//	//c.JSON(http.StatusOK, gin.H{"deleted": student})
//
//}
