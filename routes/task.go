package routes

import (
	"CanvasApplication/middleware"
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTaskRoutes(router *gin.Engine, db *gorm.DB) {
	taskController := service.NewTaskController(db)

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.POST("/", middleware.RoleMiddleware("teacher"), taskController.CreateTask)
		taskRoutes.GET("/course/:courseId", taskController.GetTasksByCourse)
	}
}
