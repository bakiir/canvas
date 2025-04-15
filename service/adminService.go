package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

//==================================GetAllUsers==================================

func (ac *AdminController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := ac.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

//==================================CreateUser==================================

func (ac *AdminController) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"login" binding:"required, email"`
		Password string `json:"password" binding:"required,min=4"`
		Role     string
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности логина
	var existing models.User
	if err := ac.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	//hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if input.Role == "teacher" || input.Role == "student" {
		user.Role = input.Role
	} else {
		user.Role = "student"
	}

	if err := ac.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create " + user.Role})
		return
	}

	// Не возвращаем пароль в ответе
	user.Password = ""
	c.JSON(
		http.StatusCreated,
		gin.H{"data": user},
	)

	c.JSON(200, gin.H{"data": user})
}

//==================================DeleteUser==================================

func (ac *AdminController) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	//  checking id's format
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalidUser ID"})
		return
	}

	//  checking user axisting
	var user models.User
	if err := ac.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	//  deleting
	if err := ac.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

//==================================GetAllCourses==================================

func (ac *AdminController) GetAllCourses(c *gin.Context) {
	var courses []models.Course
	if err := ac.db.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(200, courses)
}

//==================================CreateCourse==================================

func (ac *AdminController) CreateCourse(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Capacity    int    `json:"capacity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Capacity:    25,
	}

	if err := ac.db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course " + course.Name})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": course})
}

//==================================DeleteCourse==================================

func (ac *AdminController) DeleteCourse(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}
	var course models.Course
	if err := ac.db.First(&course, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found"})
		return
	}

	if err := ac.db.Delete(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete course" + course.Name})
		return
	}

	c.JSON(200, gin.H{"data": nil})
}

//==================================AddStudentsToCourse==================================

func (ac *AdminController) AddStudentsToCourse(c *gin.Context) {
	// Получаем courseId из URL
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	// Получаем список studentIds из тела запроса
	var studentIDs []uint
	if err := c.ShouldBindJSON(&studentIDs); err != nil {
		c.JSON(400, gin.H{"error": "Invalid student ID list"})
		return
	}

	// Находим курс
	var course models.Course
	if err := ac.db.Preload("Students").First(&course, courseID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found"})
		return
	}

	// Находим студентов
	var students []models.User
	if err := ac.db.Where("id IN ?", studentIDs).Find(&students).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch students"})
		return
	}

	// Добавляем студентов к курсу
	if err := ac.db.Model(&course).Association("Students").Append(&students); err != nil {
		c.JSON(500, gin.H{"error": "Failed to add students to course"})
		return
	}

	// Обновляем данные курса с присоединёнными студентами
	if err := ac.db.Preload("Students").First(&course, courseID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to reload course"})
		return
	}

	c.JSON(200, course)
}

//==================================RemoveStudentFromCourse==================================

func (ac *AdminController) RemoveStudentFromCourse(c *gin.Context) {
	courseIdStr := c.Param("courseId")
	studentIdStr := c.Param("studentId")

	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	studentId, err := strconv.ParseUint(studentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var course models.Course
	if err := ac.db.Preload("Students").First(&course, courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var student models.User
	if err := ac.db.First(&student, studentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if err := ac.db.Model(&course).Association("Students").Delete(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove student from course"})
		return
	}

	// Обновим курс с новым списком студентов
	if err := ac.db.Preload("Students").First(&course, courseId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload course"})
		return
	}

	c.JSON(http.StatusOK, course)
}

//==================================assignTeacherToCourse==================================

func (ac *AdminController) AssignTeacherToCourse(c *gin.Context) {
	courseIdStr := c.Param("courseId")
	teacherIdStr := c.Param("teacherId")

	courseId, err := strconv.ParseUint(courseIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherId, err := strconv.ParseUint(teacherIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var teacher models.User
	if err := ac.db.First(&teacher, teacherId).Error; err != nil {
		c.JSON(404, gin.H{"error": "Teacher not found"})
	}

	var course models.Course
	if err := ac.db.First(&course, courseId).Error; err != nil {
		c.JSON(404, gin.H{"error": "Course not found"})
	}

	if err := ac.db.Model(&course).Association("Teacher").Append(&teacher); err != nil {
		c.JSON(500, gin.H{"error": "Failed to assign teacher " + teacher.Name + " to course " + course.Name})
		return
	}

	c.JSON(200, course)

}
