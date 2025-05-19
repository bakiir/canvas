package main

import (
	"CanvasApplication/config"
	"CanvasApplication/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация Ginup
	router := gin.Default()

	// Подключение к базе данных
	config.Connect()

	// Регистрация всех маршрутов
	registerRoutes(router)

	// Запуск сервера
	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	// Группа маршрутов для студентов
	routes.SetupStudentRoutes(router, config.DB)

	// Группа маршрутов для преподавателей
	routes.SetupTeacherRoutes(router, config.DB)

	// Дополнительные маршруты (если есть)
	routes.SetupAdminRoutes(router, config.DB)

	routes.SetupTaskRoutes(router, config.DB)

	routes.SetupHomeWorkRoutes(router, config.DB)
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
