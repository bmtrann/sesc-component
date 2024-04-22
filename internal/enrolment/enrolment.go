package enrolment

import (
	"github.com/bmtrann/sesc-component/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *EnrolmentHandler, middleware *middleware.MiddlewareHandler) {
	routes := router.Group("/courses", middleware.Handle)
	{
		routes.GET("/list", handler.List)
		routes.GET("/:id/enrolment", handler.View)
		routes.POST("/enrol", handler.Enrol)
	}
}
