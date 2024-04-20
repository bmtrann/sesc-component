package auth

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, handler *AuthHandler) {
	routes := router.Group("/auth")
	{
		routes.POST("/login", handler.Login)
		routes.POST("/register", handler.Register)
	}
}
