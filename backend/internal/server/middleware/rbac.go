package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(ContextRoleKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		currentRole, ok := v.(string)
		if !ok || currentRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			return
		}

		c.Next()
	}
}
