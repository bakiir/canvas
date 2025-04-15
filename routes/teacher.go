package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTeacherRoutes(router *gin.Engine, db *gorm.DB) {
	teacherController := service.NewTeacherController(db)

	teacherRoutes := router.Group("/teachers")
	{
		teacherRoutes.GET("/courses/:teacherId", teacherController.GetCoursesForTeacher)
		teacherRoutes.POST("/courses/:courseId/assignments", teacherController.AddAssignmentToCourse)
		teacherRoutes.DELETE("/assignments/:assignmentId", teacherController.RemoveAssignmentFromCourse)
	}
}
