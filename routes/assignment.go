package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAssignmentRoutes(router *gin.Engine, db *gorm.DB) {
	AssignmentController := service.NewAssignmentController(db)
	assignmentRoutes := router.Group("/api/assignments")
	{
		assignmentRoutes.GET("/course/:courseId", AssignmentController.GetAssignmentsByCourse)

	}

}
