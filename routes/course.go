package routes

import (
	"CanvasApplication/service"
	"github.com/gin-gonic/gin"
)

func courseRoutes(router *gin.Engine) {
	router.GET("/", service.GetAllCourses)
}
