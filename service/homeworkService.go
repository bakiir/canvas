package service

import (
	"CanvasApplication/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
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
	// Получаем параметры из формы
	studentIDStr := c.PostForm("student_id")
	taskIDStr := c.PostForm("task_id")

	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student_id"})
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task_id"})
		return
	}

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open uploaded file"})
		return
	}
	defer src.Close()

	// Отправляем файл в файловый сервис (порт 8000)
	fileURL, err := sendFileToFileService(file.Filename, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to file service"})
		return
	}

	// Создаем запись в БД
	homework := models.Homework{
		StudentID:  uint(studentID),
		TaskID:     uint(taskID),
		FileURL:    fileURL,
		UploadedAt: time.Now(),
	}

	if err := ctrl.DB.Create(&homework).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Homework uploaded successfully",
		"homework": homework,
	})
}

func (ctrl *HomeworkController) GetListOfHomeworks() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskId")

		var homeworks []models.Homework
		if err := ctrl.DB.Preload("Student").Where("task_id = ?", taskID).Find(&homeworks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении домашек"})
			return
		}

		c.JSON(http.StatusOK, homeworks)
	}
}

func (ctrl *HomeworkController) DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskId")
		studentID := c.Param("studentId")

		var homework models.Homework
		if err := ctrl.DB.Where("task_id = ? AND student_id = ?", taskID, studentID).First(&homework).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Домашнее задание не найдено"})
			return
		}

		// Получаем имя файла из сохранённого URL (например, "http://localhost:8000/download?file=myfile.png")
		fileName := filepath.Base(homework.FileURL)

		// Формируем URL до файлового сервиса
		fileServiceURL := fmt.Sprintf("http://localhost:8000/download?file=%s", url.QueryEscape(fileName))

		// Делаем GET-запрос к файловому сервису
		resp, err := http.Get(fileServiceURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения файла с файлового сервиса"})
			return
		}
		defer resp.Body.Close()

		// Копируем заголовки (только нужные)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
		c.Header("Content-Type", "application/octet-stream")

		// Потоково копируем тело ответа в клиентский ответ
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке файла"})
			return
		}
	}
}

func sendFileToFileService(filename string, file io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:8000/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		URL string `json:"file_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}

func (ctrl *HomeworkController) SetGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskId")
		studentID := c.Param("studentId")

		// Только учитель может ставить оценки
		userRole := c.GetString("role") // предполагаем, что роль уже вытянута из middleware
		if userRole != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Только учитель может выставлять оценки"})
			return
		}

		// Парсим тело запроса
		var request struct {
			Grade uint `json:"grade"`
		}
		if err := c.BindJSON(&request); err != nil || request.Grade > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Оценка должна быть числом от 0 до 100"})
			return
		}

		// Находим работу
		var homework models.Homework
		if err := ctrl.DB.Where("task_id = ? AND student_id = ?", taskID, studentID).First(&homework).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Домашняя работа не найдена"})
			return
		}

		// Обновляем оценку
		homework.Grade = &request.Grade
		if err := ctrl.DB.Save(&homework).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить оценку"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Оценка успешно выставлена", "grade": request.Grade})
	}
}

func (ctrl *HomeworkController) GetStudentGrades() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем studentID и taskID из URL параметров
		taskIDStr := c.Param("taskId")
		taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		// Получаем user_id из контекста (его положил RoleMiddleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		userIDFloat, ok := userIDInterface.(float64) // jwt.MapClaims хранит числа как float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}
		userID := uint(userIDFloat)

		// Находим домашнее задание для данного студента и задания
		var homework models.Homework
		if err := ctrl.DB.Where("student_id = ? AND task_id = ?", userID, taskID).First(&homework).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Homework not found or access denied"})
			return
		}

		// Возвращаем оценку студенту
		c.JSON(http.StatusOK, gin.H{
			"task_id":    homework.TaskID,
			"student_id": homework.StudentID,
			"grade":      homework.Grade,
		})
	}
}
