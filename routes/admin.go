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
		adminRoutes.POST("/students", adminController.CreateStudent)
		adminRoutes.DELETE("/students/:id", adminController.DeleteStudent)

		// POST /admin/teachers - создание преподавателя администратором
		adminRoutes.POST("/teachers", adminController.CreateTeacher)
		adminRoutes.DELETE("/teachers/:id", adminController.DeleteTeacher)

		//adminRoutes.DELETE("/:id/courses", adminController.DeleteCourse)

	}
}
