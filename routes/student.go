package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupStudentRoutes(router *gin.Engine, db *gorm.DB) {
	studentController := service.NewStudentController(db)

	studentRoutes := router.Group("/students")
	{
		studentRoutes.GET("/", studentController.GetAllStudents)
		studentRoutes.POST("/login", studentController.Login)
		studentRoutes.POST("/enroll", studentController.Enroll)
		studentRoutes.GET("/:id/courses", studentController.GetStudentCourses)
	}
}
