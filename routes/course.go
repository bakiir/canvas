package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCourseRoutes(router *gin.Engine, db *gorm.DB) {
	courseController := service.NewCourseController(db)

	courseRoutes := router.Group("/api/courses")
	//courseRoutes.Use(middleware.AdminAuth()) // Раскомментируйте после реализации middleware
	{
		courseRoutes.GET("/teacher/:teacher_id", courseController.GetCoursesByTeacher)
		courseRoutes.GET("/student/:studentId", courseController.GetCoursesByStudent)
		courseRoutes.DELETE("/:id", courseController.DeleteCourse)
	}
}
