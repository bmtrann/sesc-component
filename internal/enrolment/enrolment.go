package enrolment

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *EnrolmentHandler) {
	routes := router.Group("/courses")
	{
		routes.GET("/list", handler.List)
		routes.GET("/:id/enrolment", handler.View)
		routes.POST("/enrol", handler.Enrol)
	}
}
