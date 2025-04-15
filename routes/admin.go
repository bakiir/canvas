package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAdminRoutes(router *gin.Engine, db *gorm.DB) {
	adminController := service.NewAdminController(db)

	// Группа маршрутов /admin с middleware для аутентификации администратора
	adminRoutes := router.Group("/admin")
	//adminRoutes.Use(middleware.AdminAuth()) // Раскомментируйте после реализации middleware
	{
		// POST /admin/students - создание студента администратором
		adminRoutes.GET("/users", adminController.GetAllUsers)
		adminRoutes.POST("/users", adminController.CreateUser)
		adminRoutes.DELETE("/users/:id", adminController.DeleteStudent)

		adminRoutes.GET("/courses", adminController.GetAllCourses)
		adminRoutes.POST("/courses", adminController.CreateCourse)
		adminRoutes.DELETE("/courses/:id", adminController.DeleteCourse)
		adminRoutes.POST("/courses/:id?students", adminController.AddStudentsToCourse)
		adminRoutes.DELETE("/courses/{courseId}/students/{studentId}", adminController.RemoveStudentFromCourse)
		adminRoutes.POST("/courses/{courseId}/teachers/{teacherId}", adminController.AssignTeacherToCourse)

	}
}
