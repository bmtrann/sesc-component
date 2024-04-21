package profile

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *ProfileHandler) {
	routes := router.Group("/profile")
	{
		routes.GET("/:id", handler.GetProfile)
		routes.GET("/:id/graduation-status", handler.GetGraduationStatus)
		routes.POST("/update", handler.UpdateProfile)
	}
}
