package routes

import (
	"CanvasApplication/middleware"
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTeacherRoutes(router *gin.Engine, db *gorm.DB) {
	teacherController := service.NewTeacherController(db)

	teacherRoutes := router.Group("/teachers")
	{
		teacherRoutes.GET("/", teacherController.GetAllTeachers)
		teacherRoutes.POST("/register", teacherController.CreateTeacher)

		teacherRoutes.POST("/login", teacherController.TeacherLogin)
		teacherRoutes.GET("/:id/courses", teacherController.GetTeacherCourses)

		teacherRoutes.POST("/courses", middleware.AuthMiddleware(), teacherController.CreateCourse)

	}
}
