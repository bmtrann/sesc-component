package profile

import (
	"github.com/bmtrann/sesc-component/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *ProfileHandler, middleware *middleware.MiddlewareHandler) {
	routes := router.Group("/profile", middleware.Handle)
	{
		routes.GET("/:id", handler.GetProfile)
		routes.GET("/:id/graduation-status", handler.GetGraduationStatus)
		routes.POST("/update", handler.UpdateProfile)
	}
}
