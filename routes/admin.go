package routes

import (
	"CanvasApplication/middleware"
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAdminRoutes(router *gin.Engine, db *gorm.DB) {
	adminController := service.NewAdminController(db)

	adminRoutes := router.Group("/admin")

	{

		adminRoutes.POST("/login", adminController.Login)
		adminRoutes.POST("/register", adminController.Register)

		adminRoutes.DELETE("/students/:id",
			middleware.RoleMiddleware("admin"),
			adminController.DeleteStudent)

		adminRoutes.DELETE("/teachers/:id",
			middleware.RoleMiddleware("admin"),
			adminController.DeleteTeacher)

		adminRoutes.DELETE("/courses/:id",
			middleware.RoleMiddleware("admin"),
			adminController.DeleteCourse)

	}
}
