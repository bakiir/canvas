package routes

import (
	"CanvasApplication/middleware"
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupHomeWorkRoutes(router *gin.Engine, db *gorm.DB) {
	homeworkController := service.NewHomeworkController(db)

	homeworkRoutes := router.Group("/homeworks")
	{
		homeworkRoutes.POST("/upload", middleware.RoleMiddleware("student"),
			homeworkController.UploadHomework)
		homeworkRoutes.GET("/:taskId/download/:studentId",
			middleware.RoleMiddleware("teacher", "student"),
			homeworkController.DownloadFile())

		homeworkRoutes.GET("/:taskId/homeworks",
			middleware.RoleMiddleware("teacher"),
			homeworkController.GetListOfHomeworks())

	}
}
