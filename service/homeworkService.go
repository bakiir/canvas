package service

import (
	"CanvasApplication/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type HomeworkController struct {
	DB *gorm.DB
}

func NewHomeworkController(db *gorm.DB) *HomeworkController {
	return &HomeworkController{DB: db}
}

type HomeworkUploadInput struct {
	StudentID uint   `json:"student_id"`
	TaskID    uint   `json:"task_id"`
	FileURL   string `json:"file_url"` // от файл-сервиса
}

func (ctrl *HomeworkController) UploadHomework(c *gin.Context) {
	var input HomeworkUploadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	homework := models.Homework{
		StudentID:  input.StudentID,
		TaskID:     input.TaskID,
		FileURL:    input.FileURL,
		UploadedAt: time.Now(),
	}

	if err := ctrl.DB.Create(&homework).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, homework)
}
