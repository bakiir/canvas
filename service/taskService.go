package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type TaskController struct {
	db *gorm.DB
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{db: db}
}

func (ctrl *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования курса
	var course models.Course
	if err := ctrl.db.First(&course, task.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id: course not found"})
		return
	}

	if err := ctrl.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) GetTasksByCourse(c *gin.Context) {
	courseID := c.Param("courseId")
	var tasks []models.Task

	if err := ctrl.db.Where("course_id = ?", courseID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
