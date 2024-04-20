package profile

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *ProfileHandler) {
	routes := router.Group("/profile")
	{
		routes.GET("/:id", handler.GetProfile)
		router.GET("/:id/graduation-status", handler.GetGraduationStatus)
		router.POST("/update", handler.UpdateProfile)
	}
}
