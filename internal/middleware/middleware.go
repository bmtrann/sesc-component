package middleware

import (
	"net/http"
	"strings"

	"github.com/bmtrann/sesc-component/internal/auth"
	"github.com/bmtrann/sesc-component/internal/exception"
	"github.com/gin-gonic/gin"
)

type MiddlewareHandler struct {
	authHanlder *auth.AuthHandler
}

func NewMiddlewareHandler(authHandler *auth.AuthHandler) *MiddlewareHandler {
	return &MiddlewareHandler{
		authHanlder: authHandler,
	}
}

func (handler *MiddlewareHandler) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := handler.authHanlder.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		if err == exception.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		c.AbortWithStatus(status)
		return
	}
	c.Set("CurrentUser", user)
}
