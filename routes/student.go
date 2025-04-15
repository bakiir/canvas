package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupStudentRoutes(router *gin.Engine, db *gorm.DB) {
	studentController := service.NewStudenController(db)

	studentRoutes := router.Group("/students")
	{
		studentRoutes.GET("/courses/:studentId", studentController.GetCoursesForStudent)
		studentRoutes.GET("/course/:courseId", studentController.GetCourseInfo)
	}
}
