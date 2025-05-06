package routes

import (
	"CanvasApplication/middleware"
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
		adminRoutes.POST("/students", middleware.RoleMiddleware("admin"), adminController.CreateStudent)
		adminRoutes.DELETE("/students/:id", middleware.RoleMiddleware("admin"), adminController.DeleteStudent)

		// POST /admin/teachers - создание преподавателя администратором
		adminRoutes.POST("/teachers", middleware.RoleMiddleware("admin"), adminController.CreateTeacher)
		adminRoutes.DELETE("/teachers/:id", middleware.RoleMiddleware("admin"), adminController.DeleteTeacher)

		//adminRoutes.DELETE("/:id/courses", adminController.DeleteCourse)

	}
}
