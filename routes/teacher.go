package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTeacherRoutes(router *gin.Engine, db *gorm.DB) {
	//teacherController := service.NewTeacherController(db)
	//
	//teacherRoutes := router.Group("/teachers")
	//{
	//	teacherRoutes.GET("/", teacherController.GetAllTeachers)
	//	teacherRoutes.POST("/login", teacherController.TeacherLogin)
	//	teacherRoutes.POST("/courses", teacherController.CreateCourse)
	//	teacherRoutes.GET("/:id/courses", teacherController.GetTeacherCourses)
	//}
}
